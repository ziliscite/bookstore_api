CREATE TABLE Addresses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255),
    city VARCHAR(255),
    postal_code VARCHAR(100),
    country VARCHAR(100)
);

CREATE TYPE PAYMENT_STATUS AS ENUM ('pending', 'completed', 'canceled');

CREATE TABLE PaymentResults (
    id SERIAL PRIMARY KEY,
    status PAYMENT_STATUS,
    email_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE PAYMENT_METHOD AS ENUM ('PayPal', 'Bank', 'QRIS');

CREATE TABLE Orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(id) ON DELETE CASCADE,
    address_id INT REFERENCES Addresses(id) ON DELETE CASCADE,

    product_price DECIMAL(10, 2),
    tax_fee DECIMAL(10, 2),
    total_price DECIMAL(10, 2),

    payment_method PAYMENT_METHOD,
    payment_result_id INT REFERENCES PaymentResults(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE OrderBooks (
    id SERIAL PRIMARY KEY,
    book_quantity INT,

    -- Aggregate total price

    book_id INT REFERENCES Books(id) ON DELETE CASCADE,
    order_id INT REFERENCES Orders(id) ON DELETE CASCADE
);
