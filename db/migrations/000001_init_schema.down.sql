-- Drop Reviews table first since it depends on Users and Books
DROP TABLE IF EXISTS Reviews;

-- Drop BookCategory table since it depends on Books and Categories
DROP TABLE IF EXISTS BookCategory;

-- Drop Categories table
DROP TABLE IF EXISTS Categories;

-- Drop OrderBooks table since it depends on Orders and Books
DROP TABLE IF EXISTS OrderBooks;

-- Drop Books table
DROP TABLE IF EXISTS Books;

-- Drop Orders table since it depends on Users, Addresses, and PaymentResults
DROP TABLE IF EXISTS Orders;

-- Drop PaymentResults table
DROP TABLE IF EXISTS PaymentResults;

-- Drop Addresses table
DROP TABLE IF EXISTS Addresses;

-- Drop Users table
DROP TABLE IF EXISTS Users;
