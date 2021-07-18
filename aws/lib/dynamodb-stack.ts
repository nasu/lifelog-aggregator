import * as cdk from '@aws-cdk/core'
import * as dynamodb from '@aws-cdk/aws-dynamodb'

export class LifeLogDynamoDBStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)
    // Structure
    //   partition_key: activiety or visit
    //   sort_key:
    //     pk=activity: userID + duration.startTimestampMs
    //     pk=visit:    userID + duration.startTimestampMs
    new dynamodb.Table(this, id + "-google-location-history", {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      tableName: 'lifelog-metrics',
      partitionKey: {
        name: 'partition_key',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'sort_key',
        type: dynamodb.AttributeType.STRING,
      },
    })

    // Structure
    //   partition_key: session
    //   sort_key:
    //     pk=session: userID
    new dynamodb.Table(this, id + "-lifelog", {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      tableName: 'lifelog',
      partitionKey: {
        name: 'partition_key',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'sort_key',
        type: dynamodb.AttributeType.STRING,
      },
    })
  }
}