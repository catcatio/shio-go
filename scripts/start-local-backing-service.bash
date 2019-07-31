#!/usr/bin/env bash
docker run -d --rm --name shio-postgres -e POSTGRES_PASSWORD=shio-local -p 5432:5432 postgres
docker run -d --rm -p "8000:8000" -p "8004:8004" -p "11626:11626" --name shio-stellar octofoxio/stellar-integration-test-network
docker run -d --rm -p "3000:3000" --name shio-stellarexplorer octofoxio/stellarexplorer

core_url=http://localhost:11626

curl_response() {
  response=$(curl --write-out %{http_code} --silent --output /dev/null $1)
  echo ${response}
}

check_core_available () {
  url=${core_url}/info
  echo $(curl_response ${url})
}

update_base_reserve() {
  url="${core_url}/upgrades?mode=set&upgradetime=1970-01-01T00:00:00Z&basereserve=5000000"
  echo $(curl_response ${url})
}

while [[ $(check_core_available) -ne 200 ]]
do
   echo "Waiting for stellar core to be available..."
   sleep 1
done

# do not have to upgrades basereserve, as it has been done in niponchi/stellar-integration-test-network
#echo update base reserve: $(update_base_reserve)
echo stellar core started
