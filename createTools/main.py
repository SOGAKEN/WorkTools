import os
import sys
import boto3
from botocore.exceptions import ClientError

def lambda_handler(event, context):
    """
    RDS Auroraインスタンスの状態を確認し、稼働中のインスタンスを停止するLambda関数

    Args:
      event: Lambdaイベントオブジェクト
      context: Lambdaコンテキストオブジェクト

    Returns:
      処理結果メッセージ
    """

    # 環境変数からインスタンスIDリストを取得
    instance_ids = os.environ.get('INSTANCE_IDS', '').split(',')

    if not instance_ids or instance_ids == ['']:
        return "環境変数INSTANCE_IDSが設定されていません。"

    # Boto3 RDSクライアントの初期化
    try:
        client = boto3.client('rds')
    except ClientError as e:
        print(f"AWSクライアントの初期化中にエラーが発生しました: {e}")
        sys.exit(1)

    # インスタンスごとに処理
    for instance_id in instance_ids:
        instance_id = instance_id.strip()  # 余分なスペースの削除

        if not instance_id:
            continue  # 空のインスタンスIDはスキップ

        try:
            response = client.describe_db_instances(DBInstanceIdentifier=instance_id)
            db_instance_status = response['DBInstances'][0]['DBInstanceStatus']

            # インスタンスが稼働中の場合に停止
            if db_instance_status == 'available':
                client.stop_db_instance(DBInstanceIdentifier=instance_id)
                print(f"インスタンス {instance_id} を停止しました。")
            else:
                print(f"インスタンス {instance_id} は '{db_instance_status}' 状態です。操作は行いません。")

        except ClientError as e:
            print(f"インスタンス {instance_id} の処理中にエラーが発生しました: {e}")

    return "操作が完了しました。"
