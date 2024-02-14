import boto3
from datetime import datetime, timedelta

# AWS SDK (boto3)のクライアントを初期化
rds_client = boto3.client('rds')
region = 'ap-northeast-1'  # 東京リージョン

def lambda_handler(event, context):
    utc_now = datetime.utcnow()
    jst_now = utc_now + timedelta(hours=9)  # UTCからJSTへ変換

    # RDSインスタンスの処理
    instances = rds_client.describe_db_instances()
    for instance in instances['DBInstances']:
        tags = rds_client.list_tags_for_resource(ResourceName=instance['DBInstanceArn'])['TagList']
        start_time, end_time = get_time_from_tags(tags)
        if start_time and end_time:
            process_instance(instance, start_time, end_time, jst_now)

    # Auroraクラスターの処理
    clusters = rds_client.describe_db_clusters()
    for cluster in clusters['DBClusters']:
        tags = rds_client.list_tags_for_resource(ResourceName=cluster['DBClusterArn'])['TagList']
        start_time, end_time = get_time_from_tags(tags)
        if start_time and end_time:
            process_cluster(cluster, start_time, end_time, jst_now)

def get_time_from_tags(tags):
    start_time = end_time = None
    for tag in tags:
        if tag['Key'] == 'StartTime':
            start_time = tag['Value']
        elif tag['Key'] == 'EndTime':
            end_time = tag['Value']
    return start_time, end_time

def process_instance(instance, start_time_str, end_time_str, current_time):
    start_time = datetime.strptime(str(current_time.date()) + ' ' + start_time_str, '%Y-%m-%d %H:%M')
    end_time = datetime.strptime(str(current_time.date()) + ' ' + end_time_str, '%Y-%m-%d %H:%M')
    # 日付をまたぐ設定の対応
    if current_time.hour < int(start_time_str.split(':')[0]):
        start_time -= timedelta(days=1)
        end_time -= timedelta(days=1)
    if start_time <= current_time <= end_time and instance['DBInstanceStatus'] == 'stopped':
        print(f"Starting RDS instance: {instance['DBInstanceIdentifier']}")
        rds_client.start_db_instance(DBInstanceIdentifier=instance['DBInstanceIdentifier'])
    elif not(start_time <= current_time <= end_time) and instance['DBInstanceStatus'] == 'available':
        print(f"Stopping RDS instance: {instance['DBInstanceIdentifier']}")
        rds_client.stop_db_instance(DBInstanceIdentifier=instance['DBInstanceIdentifier'])

def process_cluster(cluster, start_time_str, end_time_str, current_time):
    start_time = datetime.strptime(str(current_time.date()) + ' ' + start_time_str, '%Y-%m-%d %H:%M')
    end_time = datetime.strptime(str(current_time.date()) + ' ' + end_time_str, '%Y-%m-%d %H:%M')
    # 日付をまたぐ設定の対応
    if current_time.hour < int(start_time_str.split(':')[0]):
        start_time -= timedelta(days=1)
        end_time -= timedelta(days=1)
    if start_time <= current_time <= end_time and cluster['Status'] == 'stopped':
        print(f"Starting Aurora cluster: {cluster['DBClusterIdentifier']}")
        rds_client.start_db_cluster(DBClusterIdentifier=cluster['DBClusterIdentifier'])
    elif not(start_time <= current_time <= end_time) and cluster['Status'] == 'available':
        print(f"Stopping Aurora cluster: {cluster['DBClusterIdentifier']}")
        rds_client.stop_db_cluster(DBClusterIdentifier=cluster['DBClusterIdentifier'])

if __name__ == "__main__":
    lambda_handler(None, None)
