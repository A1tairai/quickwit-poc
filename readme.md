cd ~/quickwit

aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 767398036361.dkr.ecr.ap-northeast-1.amazonaws.com

docker build --network=host -t quickwit-poc .

docker tag quickwit-poc:latest 767398036361.dkr.ecr.ap-northeast-1.amazonaws.com/quickwit-poc:latest

docker push 767398036361.dkr.ecr.ap-northeast-1.amazonaws.com/quickwit-poc:latest

kubectl delete -f deployment.yaml 

kubectl apply -f deployment.yaml 


docker run -e SEND_TO_KAFKA="true" -e KAFKA_BROKER="b-8.quickwitmsk.8otisf.c4.kafka.ap-northeast-1.amazonaws.com:9092,b-9.quickwitmsk.8otisf.c4.kafka.ap-northeast-1.amazonaws.com:9092,b-5.quickwitmsk.8otisf.c4.kafka.ap-northeast-1.amazonaws.com:9092" -e KAFKA_TOPIC="quickwig_topic" 767398036361.dkr.ecr.ap-northeast-1.amazonaws.com/quickwit-poc:latest


