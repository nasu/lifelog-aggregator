#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { LifeLogDynamoDBStack } from '../lib/dynamodb-stack';

const app = new cdk.App();

new LifeLogDynamoDBStack(app, "LifeLogDynamoDBStack")