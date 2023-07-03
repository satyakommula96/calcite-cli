/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */


package keywords

import "github.com/c-bata/go-prompt"

func CustomCompleter(d prompt.Document) []prompt.Suggest {
    input := d.GetWordBeforeCursor()
    if input == "" {
        return nil // Return an empty suggestion list when the input is empty
    }

    suggestions := []prompt.Suggest{
		{Text: "SELECT", Description: "Retrieve data from a table"},
		{Text: "FROM", Description: "Specify the table or tables to retrieve data from"},
		{Text: "WHERE", Description: "Filter rows based on a condition"},
		{Text: "JOIN", Description: "Combine rows from multiple tables based on a related column between them"},
		{Text: "GROUP BY", Description: "Group rows based on a specified column"},
		{Text: "ORDER BY", Description: "Sort rows based on one or more columns"},
		{Text: "INSERT INTO", Description: "Insert data into a table"},
		{Text: "UPDATE", Description: "Modify data in a table"},
		{Text: "DELETE FROM", Description: "Delete rows from a table"},
		{Text: "CREATE TABLE", Description: "Create a new table"},
		{Text: "ALTER TABLE", Description: "Modify an existing table"},
		{Text: "DROP TABLE", Description: "Delete an existing table"},
		{Text: "COUNT", Description: "Return the number of rows that match a specified condition"},
		{Text: "SUM", Description: "Calculate the sum of values in a column"},
		{Text: "AVG", Description: "Calculate the average of values in a column"},
		{Text: "MAX", Description: "Find the maximum value in a column"},
		{Text: "MIN", Description: "Find the minimum value in a column"},
		{Text: "DISTINCT", Description: "Return unique values in a column"},
		{Text: "AS", Description: "Rename a column or table"},
		{Text: "AND", Description: "Combine multiple conditions in a WHERE clause"},
		{Text: "OR", Description: "Specify alternative conditions in a WHERE clause"},
		{Text: "NOT", Description: "Negate a condition in a WHERE clause"},
		{Text: "BETWEEN", Description: "Specify a range of values"},
		{Text: "LIKE", Description: "Search for a pattern in a column"},
		{Text: "IN", Description: "Check if a value exists in a list"},
		{Text: "NULL", Description: "Represents a missing or unknown value"},
		{Text: "IS NULL", Description: "Check if a value is NULL"},
		{Text: "IS NOT NULL", Description: "Check if a value is not NULL"},
		{Text: "CASE", Description: "Perform conditional logic"},
		{Text: "WHEN", Description: "Specify a condition in a CASE statement"},
		{Text: "THEN", Description: "Define the result for a specific condition in a CASE statement"},
		{Text: "ELSE", Description: "Define the result for all other conditions in a CASE statement"},
		{Text: "END", Description: "End a CASE statement"},
		{Text: "INNER JOIN", Description: "Return rows that have matching values in both tables"},
		{Text: "LEFT JOIN", Description: "Return all rows from the left table and the matching rows from the right table"},
		{Text: "RIGHT JOIN", Description: "Return all rows from the right table and the matching rows from the left table"},
		{Text: "FULL JOIN", Description: "Return all rows when there is a match in either the left or right table"},
		{Text: "UNION", Description: "Combine the result of multiple SELECT statements"},
		{Text: "INTERSECT", Description: "Return the common rows between the result of multiple SELECT statements"},
		{Text: "EXCEPT", Description: "Return the rows from the first SELECT statement that are not in the result of the second SELECT statement"},
		{Text: "HAVING", Description: "Filter groups based on a condition in a GROUP BY clause"},
		{Text: "LIMIT", Description: "Limit the number of rows returned by a query"},
		{Text: "OFFSET", Description: "Skip a specified number of rows before starting to return rows"},
		{Text: "TOP", Description: "Return the top n rows from a query result"},
		{Text: "CASCADE", Description: "Automatically propagate changes to related tables"},
		{Text: "PRIMARY KEY", Description: "Define a column as the primary key for a table"},
		{Text: "FOREIGN KEY", Description: "Define a column as a foreign key to establish a relationship with another table"},
		{Text: "INDEX", Description: "Create an index on one or more columns for faster data retrieval"},
		{Text: "UNIQUE", Description: "Enforce uniqueness of values in a column or a group of columns"},
		{Text: "CHECK", Description: "Specify a condition that must be true for a row to be valid"},
		{Text: "DEFAULT", Description: "Specify a default value for a column"},
		{Text: "NULLIF", Description: "Return NULL if two expressions are equal"},
		{Text: "COALESCE", Description: "Return the first non-null expression in a list"},
		{Text: "EXISTS", Description: "Check if a subquery returns any rows"},
		{Text: "ALL", Description: "Return all rows, including duplicates"},
		{Text: "ANY", Description: "Check if any of the subquery values meet a condition"},
		{Text: "ASC", Description: "Sort rows in ascending order"},
		{Text: "DESC", Description: "Sort rows in descending order"},
		{Text: "OUTER JOIN", Description: "Return all rows from both tables, including unmatched rows"},
		{Text: "OVER", Description: "Perform a calculation across rows"},
		{Text: "PARTITION BY", Description: "Divide the rows into groups for calculation"},
		{Text: "SET", Description: "Assign values to columns in an UPDATE statement"},
		{Text: "FETCH", Description: "Retrieve a specific number of rows from a query result"},
		{Text: "FOR", Description: "Specify the locking behavior for a query"},
		{Text: "MINUS", Description: "Return the rows from the first SELECT statement that are not in the result of the second SELECT statement"},
		{Text: "ON", Description: "Specify the join condition between tables"},
		{Text: "VALUES", Description: "Specify the values to be inserted into a table"},
		{Text: "DISTINCT", Description: "Return unique values in a column"},
	}
	

    return prompt.FilterHasPrefix(suggestions, input, true)
}

