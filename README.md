
go run github.com/project-azorian/clementine/main.go
CGO_ENABLED=0 GOOS=linux go build -o ${GOPATH}/bin/clementine -a -tags netgo -ldflags '-w' github.com/project-azorian/clementine/main.go

(
./tools/deployment/developer/nfs/030-ingress.sh
./tools/deployment/developer/nfs/040-nfs-provisioner.sh ;
#NOTE: Lint and package chart
: ${OSH_INFRA_PATH:="../openstack-helm-infra"}
make -C ${OSH_INFRA_PATH} mariadb
helm upgrade --install mariadb ${OSH_INFRA_PATH}/mariadb \
    --namespace=openstack
./tools/deployment/common/wait-for-pods.sh openstack
)

DB_NAME=clementine
DB_USER=clementine
DB_PASSWORD=password
mysql \
  --host='mariadb.openstack.svc.cluster.local' \
  --port='3306' \
  --user='root' \
  --password='password' \
  --execute="\
      CREATE DATABASE IF NOT EXISTS $DB_NAME ; \
      CREATE USER IF NOT EXISTS '$DB_USER'@'%' IDENTIFIED BY '$DB_PASSWORD' ; \
      GRANT ALL ON $DB_NAME.* TO '$DB_USER'@'%' ; \
      FLUSH PRIVILEGES ;"

mysql \
  --host='mariadb.openstack.svc.cluster.local' \
  --port='3306' \
  --user="$DB_USER" \
  --password="$DB_PASSWORD" \
  --execute="\
      USE $DB_NAME ; \
      SHOW TABLES ;"

TABLE_NAME=nodes

mysql \
  --host='mariadb.openstack.svc.cluster.local' \
  --port='3306' \
  --user="$DB_USER" \
  --password="$DB_PASSWORD" \
  --execute="\
      USE $DB_NAME ; \
      create table $TABLE_NAME(
         id INT NOT NULL AUTO_INCREMENT,
         name VARCHAR(100) NOT NULL,
         PRIMARY KEY ( id )
      );"

mysql \
  --host='mariadb.openstack.svc.cluster.local' \
  --port='3306' \
  --user="$DB_USER" \
  --password="$DB_PASSWORD" \
  --execute="\
      USE $DB_NAME ; \
      insert into $TABLE_NAME(name) values('db node one') ; \
      insert into $TABLE_NAME(name) values('db node two') ;"

mysql \
  --host='mariadb.openstack.svc.cluster.local' \
  --port='3306' \
  --user="$DB_USER" \
  --password="$DB_PASSWORD" \
  --execute="\
      USE $DB_NAME ; \
      select * from $TABLE_NAME;"

curl --verbose -X GET http://localhost:8000/nodes | jq -r

curl --verbose -X GET http://localhost:8000/nodes/123 | jq -r

tee /tmp/add-node.json <<EOF
{"id":2,"name":"db node twelvty"}
EOF
curl --verbose -X POST --data "@/tmp/add-node.json" http://localhost:8000/nodes

tee /tmp/add-node.json <<EOF
{"id":114,"name":"db twelvddty"}
EOF
curl --verbose -X PUT --data "@/tmp/add-node.json" http://localhost:8000/nodes

curl --verbose -X DELETE http://localhost:8000/nodes/123 | jq -r
