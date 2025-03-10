# vinyl store





command to get local PostGres DB up and going.
```bash
docker run --name vinyl-store-storage -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
```