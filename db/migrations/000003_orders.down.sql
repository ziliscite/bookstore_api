-- Drop OrderBooks table since it depends on Orders and Books
DROP TABLE IF EXISTS OrderBooks;

-- Drop Orders table since it depends on Users, Addresses, and PaymentResults
DROP TABLE IF EXISTS Orders;

-- Drop PaymentResults table
DROP TABLE IF EXISTS PaymentResults;

-- Drop Addresses table
DROP TABLE IF EXISTS Addresses;