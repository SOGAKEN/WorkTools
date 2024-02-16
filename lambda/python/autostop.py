import boto3

def lambda_handler(event, context):
    rds = boto3.client('rds')
    
    # RDSインスタンスの処理
    instances = rds.describe_db_instances()
    for instance in instances['DBInstances']:
        instance_id = instance['DBInstanceIdentifier']
        instance_status = instance['DBInstanceStatus']
        
        # インスタンスのタグを取得
        tags = rds.list_tags_for_resource(ResourceName=instance['DBInstanceArn'])['TagList']
        auto_stop_tag = next((tag for tag in tags if tag['Key'] == 'AutoStop'), None)
        
        # AutoStopタグが存在し、インスタンスが稼働中の場合に停止
        if auto_stop_tag and instance_status == 'available':
            print(f"Stopping RDS instance {instance_id}...")
            rds.stop_db_instance(DBInstanceIdentifier=instance_id)
        else:
            print(f"Skipping RDS instance {instance_id}...")

    # Auroraクラスターの処理
    clusters = rds.describe_db_clusters()
    for cluster in clusters['DBClusters']:
        cluster_id = cluster['DBClusterIdentifier']
        cluster_status = cluster['Status']
        
        # クラスターのタグを取得
        tags = rds.list_tags_for_resource(ResourceName=cluster['DBClusterArn'])['TagList']
        auto_stop_tag = next((tag for tag in tags if tag['Key'] == 'AutoStop'), None)
        
        # AutoStopタグが存在し、クラスターが稼働中の場合に停止
        if auto_stop_tag and cluster_status == 'available':
            print(f"Stopping Aurora cluster {cluster_id}...")
            rds.stop_db_cluster(DBClusterIdentifier=cluster_id)
        else:
            print(f"Skipping Aurora cluster {cluster_id}...")

if __name__ == "__main__":
    lambda_handler(None, None)
