CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,

    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,

    is_admin BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE Addresses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255),
    city VARCHAR(255),
    postal_code VARCHAR(100),
    country VARCHAR(100)
);

CREATE TABLE PaymentResults (
    id SERIAL PRIMARY KEY,
    status VARCHAR(100),
    email_address VARCHAR(255),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE Orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(id) ON DELETE CASCADE,
    address_id INT REFERENCES Addresses(id) ON DELETE CASCADE,

    product_price DECIMAL(10, 2),
    tax_fee DECIMAL(10, 2),
    total_price DECIMAL(10, 2),

    payment_method VARCHAR(100),
    payment_result_id INT REFERENCES PaymentResults(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE Books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,

    slug VARCHAR(255) UNIQUE,
    cover_image VARCHAR(255),
    synopsis TEXT,

    -- rating int // can be aggregated
    -- num_reviews // can also be aggregated

    price DECIMAL(10, 2),
    stock INT,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE OrderBooks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    quantity INT,
    price INT,

    order_id INT REFERENCES Orders(id) ON DELETE CASCADE,
    book_id INT REFERENCES Books(id) ON DELETE CASCADE
);

CREATE TABLE Categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE BookCategory (
    book_id INT REFERENCES Books(id) ON DELETE CASCADE,
    category_id INT REFERENCES Categories(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);

CREATE TABLE Reviews (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    rating INT,
    comment VARCHAR(500),

    user_id INT REFERENCES Users(id) ON DELETE CASCADE,
    book_id INT REFERENCES Books(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
