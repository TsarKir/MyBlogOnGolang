# Exercise: Blog Implementation with Rate Limiting

## Requirements

### Features

1. **Database**: 
   - Use PostgreSQL for data storage.

2. **Admin Panel**: 
   - Accessible at the `/admin` endpoint.
   - Login form for posting new articles (no editing functionality required).

3. **Pagination**: 
   - Display a maximum of 3 articles per page to avoid overwhelming users.

4. **Application UI**: 
   - The application should run on port 8888.

5. **File Structure**: 
   - Submit all additional files (images, CSS, JS) in a zip file.
   - The resulting directory structure should be as follows:
     ```
     .
     ├── css
     │   └── main.css
     ├── images
     │   └── my_cat.png
     ├── js
     │   └── scripts.js
     └── myblog-binary
     ```

6. **Admin Credentials**: 
   - Create a file named `admin_credentials.txt` containing:
     - Admin login and password for posting access.
     - Database credentials (database name and user).
     - Any additional commands needed to create tables in the database.

7. **Main Page**:
   - Include a logo.
   - Provide links to articles with optional short previews of their content.
   - Implement pagination if there are more than 3 articles.

8. **Article Page**:
   - Render the content of each article.
   - Include a "Back" link to return to the main page.

9. **Rate Limiting**:
    - Implement rate limiting to restrict access to no more than 100 clients per second.
    - If this limit is exceeded, return a "429 Too Many Requests" response.

10. **Scalability**:
    - The rate limit can be adjusted in the future as traffic increases.
