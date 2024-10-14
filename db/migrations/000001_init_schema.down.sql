-- Drop Reviews table first since it depends on Users and Books
DROP TABLE IF EXISTS Reviews;

-- Drop BookCategory table since it depends on Books and Categories
DROP TABLE IF EXISTS BookCategory;

-- Drop Categories table
DROP TABLE IF EXISTS Categories;

-- Drop Books table
DROP TABLE IF EXISTS Books;

-- Drop Users table
DROP TABLE IF EXISTS Users;
