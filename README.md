# sensor_info_service
sensor info service

service listens to nats server for sensor location change events and persists them for times when sensors are moved or lose power
this service also answers get location request
scenario - sensor loses power ( or is physically moved ), when the sensor powers up again it will ask the service for its location
currently this service has a fixed id of 000 but eventually it will be discovered uing ecomms

