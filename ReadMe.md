curl -X POST http://localhost:6080/set/peakshaving/minimalPower -H "Content-Type: application/json" -d '{"value": -2.5}'

curl http://localhost:6080/get/peakshaving/minimalPower

curl -X POST http://localhost:6080/set/peakshaving/maximalPower -H "Content-Type: application/json" -d '{"value": 50}'

curl http://localhost:6080/get/peakshaving/maximalPower

curl https://ess-server.fly.dev/get/peakshaving/maximalPower


# ess-server
