# RISTEK TEST - Cleo Excellen Iskandar✍🏻
## Pre-Requisite 🔌
Before cloning this repository make sure:
1. [Go](https://go.dev/doc/install) is installed
2. [Postgresql](https://www.postgresql.org/download/) is downloaded

## Getting Started ⏰
1. Clone the repo
```bash
git clone https://github.com/cleoexcel/ristek-test.git
cd ristek-test
```

2. Open Postgresql and create schema `ristektest`
3. Copy the content in `.env.sample` to a new `.env` file (Note: adjustment on the content may be needed)
4. Install dependencies
```bash
go mod tidy
```
5. run the project
```bash
go run main.go
```

## API Contract ⚖️
Anda dapat mengakses penjelasan secara detail pada setiap routes melalui [link ini (swagger)](https://app.swaggerhub.com/apis/CleoExcellen/OPREC_RISTEK/1.0.0) atau [link ini (postman)](https://api-ristek.postman.co/workspace/API-RISTEK-Workspace~c1cb0d07-3c89-45f6-82b9-3202ba8d08b1/collection/38268031-b06bb323-8b25-42d0-8c4d-d884686e69b2?action=share&creator=38268031)

## Penjelasan 📜
- Pada category tryout hanya terdiri dari Biologi, Physics, Math, Chemistry, dan History.
- Pengguna dapat mengedit title dan description tryout, tetapi tidak dapat mengedit category
- Pengguna dapat menjawab/mensubmit tryout dan mendapatkan score berdasarkan bobot yang diberikan pada setiap question

## Alur 👣
1. User mendaftar akunnya dengan memberikan username dan password. 
2. User melakukan login ke akun mereka.
3. User dapat membuat tryout mereka.
4. User dapat melihat tryout-tryout dari user lain dan tryout mereka sendiri
5. User dapat mengedit judul dan description tryout mereka sendiri.
6. User dapat membuat/mengedit/mendelete question sebelum ada yang melakukan submisi pada tryout mereka
7. User dapat melakukan banyak submisi pada satu tryout
8. User akan mendapatkan score pada setiap submisi yang mereka lakukan
