# Kombinat

:see_no_evil: A first try to curate [Auto Scaling](https://aws.amazon.com/de/autoscaling/) with [CoreOS](https://coreos.com/) in [AWS](https://aws.amazon.com/).

We will curate many things in the months to come, but right now, there are these things.

* [etcd](https://coreos.com/etcd/docs/latest/)
* EC2 Peers

## Features

### etcd

We curate the etcd cluster members, but removing members that are unreachable, because related AutoScaling instances have vanished.

> Please, keep in mind the (N+1)/2 limit of scaling members of an etcd cluster

## Get Started

### Setup

Setting up the needed deps and restore the relevant dependencies.

```
make deps && make restore
```

Build the project and see view the produced binaries.

```
make build && ls -l ./bin
```

> There is also a `Dockerfile` included, which allows to run `kombinat` in a container with [Supervisord](http://supervisord.org/)

> The docker has to be on the host net `docker run --rm --network host -d pixelmilk/kombinat`

### EC2 IAM Policy

> We may mention the [ec2-metadata](https://github.com/axelspringer/ec2-metadata) project, which allows to test the Instance Metadata and User Data Service :grin:

Beside the need to access the to [Instance Metadata and User Data Service](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) we need an IAM Role with the following IAM Policy for the EC2 Instances in the AutoScaling Groups.

```
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Effect": "Allow",
          "Action": "ec2:Describe*",
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": "elasticloadbalancing:Describe*",
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "cloudwatch:ListMetrics",
              "cloudwatch:GetMetricStatistics",
              "cloudwatch:Describe*"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": "autoscaling:Describe*",
          "Resource": "*"
      }
  ]
}
```

This policy allows to read the state of the AutoScaling Groups and the related EC2 instances.

# License
[MIT](/LICENSE)
