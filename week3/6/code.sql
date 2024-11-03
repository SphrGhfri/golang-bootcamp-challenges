CREATE TABLE passengers (
    passengers_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    passenger_name VARCHAR(255) NOT NULL,
    passenger_phone VARCHAR(15) NOT NULL UNIQUE
);

CREATE TABLE vehicle_types (
    vehicle_types_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    vehicle_type VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE drivers_vehicles (
    drivers_vehicles_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    driver_name VARCHAR(255) NOT NULL,
    driver_phone VARCHAR(15) NOT NULL UNIQUE,
    car_model VARCHAR(255) NOT NULL,
    car_plate VARCHAR(10) NOT NULL UNIQUE,
    vehicle_types_id BIGINT NOT NULL,
    FOREIGN KEY (vehicle_types_id) REFERENCES vehicle_types(vehicle_types_id)
);

CREATE TABLE traffic_conditions (
    traffic_conditions_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    traffic_condition VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE discounts (
    discounts_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    discount_code VARCHAR(10) NOT NULL UNIQUE
);

CREATE TABLE rides (
    rides_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    passengers_id BIGINT NOT NULL,
    drivers_vehicles_id BIGINT NOT NULL,
    pickup_location TEXT NOT NULL,
    dropoff_location TEXT NOT NULL,
    fare BIGINT NOT NULL,
    ride_date DATE NOT NULL,
    discounts_id BIGINT NOT NULL,
    traffic_conditions_id BIGINT NOT NULL,
    FOREIGN KEY (passengers_id) REFERENCES passengers(passengers_id),
    FOREIGN KEY (drivers_vehicles_id) REFERENCES drivers_vehicles(drivers_vehicles_id),
    FOREIGN KEY (discounts_id) REFERENCES discounts(discounts_id),
    FOREIGN KEY (traffic_conditions_id) REFERENCES traffic_conditions(traffic_conditions_id)
);
