#!/usr/bin/env bash
DEPLOY_ENV=dev

cloudFormationDelete()
{
    STACK_ROLLBACK=$(AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" aws cloudformation list-stacks --region "$AWS_REGION" --stack-status-filter ROLLBACK_COMPLETE UPDATE_ROLLBACK_COMPLETE | jq '.StackSummaries[].StackName//empty' | grep "$STACK_NAME")
    if [[ -z "$STACK_ROLLBACK" ]] || [[ "$STACK_ROLLBACK" == "" ]]; then
        echo ""$STACK_NAME" in good state"
    else
        echo "Deleting Stack "$STACK_NAME""
        AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY aws cloudformation delete-stack \
            --region "$AWS_REGION" \
            --stack-name "$STACK_NAME"
        sleep 30
        echo "Deleted Stack "$STACK_NAME""
    fi
}

cloudFormation()
{
    STACK_EXISTS=$(AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" aws cloudformation list-stacks --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE --region "$AWS_REGION" | jq '.StackSummaries[].StackName//empty' | grep "$STACK_NAME")
    if [[ -z "$STACK_EXISTS" ]] || [[ "$STACK_EXISTS" == "" ]]; then
        AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY aws cloudformation create-stack \
            --template-url https://"$S3_FOLDER".s3."$AWS_REGION".amazonaws.com/"$SERVICE_NAME"/cf.yaml \
            --stack-name "$STACK_NAME" \
            --region "$AWS_REGION" \
            --capabilities CAPABILITY_NAMED_IAM \
            --parameters \
                ParameterKey=ServiceName,ParameterValue=$SERVICE_NAME \
                ParameterKey=BuildKey,ParameterValue=$SERVICE_NAME/"$TRAVIS_BUILD_ID".zip \
                ParameterKey=Environment,ParameterValue="$DEPLOY_ENV"  \
                ParameterKey=BuildBucket,ParameterValue="$S3_FOLDER" \
                ParameterKey=DomainName,ParameterValue="$SERVICE_NAME"."$DNS_ZONE_NAME" \
                ParameterKey=AuthorizerARN,ParameterValue="$AUTHORIZER_ARN" \
                ParameterKey=CertificateARN,ParameterValue="$CERTIFICATE_ARN" \
                ParameterKey=DNSZoneName,ParameterValue="$DNS_ZONE_NAME".
    else
        AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY aws cloudformation update-stack \
            --template-url https://"$S3_FOLDER".s3."$AWS_REGION".amazonaws.com/"$SERVICE_NAME"/cf.yaml \
            --stack-name "$STACK_NAME" \
            --region "$AWS_REGION" \
            --capabilities CAPABILITY_NAMED_IAM \
            --parameters \
                ParameterKey=ServiceName,ParameterValue=$SERVICE_NAME \
                ParameterKey=BuildKey,ParameterValue=$SERVICE_NAME/"$TRAVIS_BUILD_ID".zip \
                ParameterKey=Environment,ParameterValue="$DEPLOY_ENV" \
                ParameterKey=BuildBucket,ParameterValue="$S3_FOLDER" \
                ParameterKey=DomainName,ParameterValue="$SERVICE_NAME"."$DNS_ZONE_NAME" \
                ParameterKey=AuthorizerARN,ParameterValue="$AUTHORIZER_ARN" \
                ParameterKey=CertificateARN,ParameterValue="$CERTIFICATE_ARN" \
                ParameterKey=DNSZoneName,ParameterValue="$DNS_ZONE_NAME".
    fi
}

deployIt()
{
    AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY aws s3 cp cf.yaml s3://"$S3_FOLDER"/$SERVICE_NAME/cf.yaml
    AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY aws s3 cp "$TRAVIS_BUILD_ID".zip s3://$S3_FOLDER/$SERVICE_NAME/"$TRAVIS_BUILD_ID".zip
}

if [[ -z "$TRAVIS_PULL_REQUEST" ]] || [[ "$TRAVIS_PULL_REQUEST" == "false" ]]; then
    AWS_ACCESS_KEY_ID=$DEV_AWS_ACCESS_KEY_ID
    AWS_SECRET_ACCESS_KEY=$DEV_AWS_SECRET_ACCESS_KEY
    S3_FOLDER=$DEV_S3_BUCKET
    DNS_ZONE_NAME=$DEV_DNS_ZONE_NAME

    echo "Deploy Dev"
    deployIt
    cloudFormationDelete
    cloudFormation
    echo "Deployed Dev"

    # Master has an extra step to launch into live
    if [[ "$TRAVIS_BRANCH" == "master" ]]; then
        AWS_ACCESS_KEY_ID=$LIVE_AWS_ACCESS_KEY_ID
        AWS_SECRET_ACCESS_KEY=$LIVE_AWS_SECRET_ACCESS_KEY
        S3_FOLDER=$LIVE_S3_BUCKET
        DNS_ZONE_NAME=$LIVE_DNS_ZONE_NAME
        DEPLOY_ENV=live

        echo "Deploy Live"
        deployIt
        cloudFormationDelete
        cloudFormation
        echo "Deployed Live"
    fi
fi