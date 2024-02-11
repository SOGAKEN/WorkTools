import boto3
from datetime import datetime, timedelta

# AWSの認証情報を設定（環境変数または~/.aws/credentialsを使用）
# RDSクライアントの初期化
rds = boto3.client('rds')

def jst_to_utc(jst_time_str):
    """JSTの時刻をUTCに変換する"""
    # JSTをdatetimeオブジェクトに変換（日付は任意でOK、ここでは現在の日付を使用）
    jst_time = datetime.strptime(jst_time_str, '%H:%M')
    # JSTからUTCに変換（9時間引く）
    utc_time = jst_time - timedelta(hours=9)
    return utc_time.strftime('%H:%M')

def tag_based_management():
    # 現在のUTC時刻を取得
    now = datetime.utcnow()
    current_time = now.strftime('%H:%M')

    # RDSインスタンスのリストを取得
    instances = rds.describe_db_instances()
    for instance in instances['DBInstances']:
        manage_instance(instance, current_time)

    # Aurora DBクラスターのリストを取得
    clusters = rds.describe_db_clusters()
    for cluster in clusters['DBClusters']:
        manage_instance(cluster, current_time, is_cluster=True)

def manage_instance(instance, current_time, is_cluster=False):
    # インスタンスまたはクラスターのARNを取得
    arn = instance['DBInstanceArn'] if not is_cluster else instance['DBClusterArn']
    tags = rds.list_tags_for_resource(ResourceName=arn)

    start_time = None
    stop_time = None
    for tag in tags['TagList']:
        if tag['Key'] == 'StartTime':
            start_time = jst_to_utc(tag['Value'])
        elif tag['Key'] == 'StopTime':
            stop_time = jst_to_utc(tag['Value'])

    # タグが設定されていない場合は何もしない
    if not start_time or not stop_time:
        return

    # 現在時刻が停止時刻未満かどうかを確認し、条件に応じてインスタンスまたはクラスターを管理
    if current_time < stop_time:
        if instance['DBInstanceStatus'] != 'available' and not is_cluster:
            # RDSインスタンスを起動
            print(f"Starting RDS instance: {instance['DBInstanceIdentifier']}")
            rds.start_db_instance(DBInstanceIdentifier=instance['DBInstanceIdentifier'])
        elif is_cluster and instance['Status'] != 'available':
            # Auroraクラスターを起動
            print(f"Starting Aurora cluster: {instance['DBClusterIdentifier']}")
            rds.start_db_cluster(DBClusterIdentifier=instance['DBClusterIdentifier'])

if __name__ == "__main__":
    tag_based_management()
