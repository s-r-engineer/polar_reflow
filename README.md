# Polar reflow

Current state: **alpha / POC**

## Abstract
The idea is to provide historic HRV data

## Architecture
Consists of:
1. Polar reflow
2. InfluxDB
3. Grafana

All is running in docker compose. **Polar reflow** built from sources

## Users manual
1. Download your polar data from https://account.polar.com/
2. Spin up docker compose in **deploy** folder with `docker compose up -d`
3. Upload data using **/uploaddata** endpoint. You can use curl like this: `curl -XPUT localhost:6969/uploaddata -F 'file=@./data.tgz'` **NB!** You can upload only archived data for now. You can reupload it safely. Data that is already in the database will be discarded and will not affect calculations 
4. login to Grafana. User / password is **admin / admin2** 
5. Only one dashboard is available, and you are able to see your HRV data there by setting timerange. 