# Golang Training Notes

## 1. Struktur Dasar Program Go

Program Go biasanya dimulai dari `package main`, lalu import package yang dibutuhkan, dan entry point program ada di `func main()`.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

Penjelasan singkat:

- `package main` berarti file ini bagian dari program yang bisa langsung dijalankan.
- `import "fmt"` dipakai untuk menggunakan package `fmt`.
- `func main()` adalah fungsi utama yang pertama kali dijalankan.
- `fmt.Println()` digunakan untuk print ke terminal.

---

## 2. Const

`const` digunakan untuk nilai yang tidak berubah.

Format dasar:

```go
const <nama> <tipe> = <value>
```

Contoh:

```go
const Gender string = "Male"
```

Bisa juga dibuat dalam bentuk grouping:

```go
const (
    GenderMale   string = "Male"
    GenderFemale string = "Female"
)
```

---

## 3. Array dan Slice

Array ukurannya tetap, sedangkan slice lebih fleksibel/dinamis.

```go
var numbers = []int{1, 2, 3, 4, 5} // slice, dynamic
var fixedNumbers = [4]int{1, 2, 3, 4} // array, fixed size
```

Kalau ingin struktur data yang bisa bertambah atau dikurangi, biasanya pakai slice.

---

## 4. Map

Map menyimpan data dalam bentuk key-value.

```go
var numbers = map[int]string{
    1: "satu",
    2: "dua",
    3: "tiga",
}
```

Cara akses:

```go
fmt.Println(numbers[1]) // output: satu
```

Cara looping map:

```go
for k, v := range numbers {
    fmt.Println(k, v)
}
```

Catatan penting:

Urutan iterasi map di Go tidak dijamin selalu sama. Misalnya map punya key `1, 2, 3`, output saat di-loop bisa saja `1 2 3`, tapi bisa juga `1 3 2`.

Kalau butuh urutan yang pasti, gunakan slice.

---

## 5. Struct

Struct digunakan untuk membuat tipe data sendiri yang punya beberapa field.

```go
type Person struct {
    Name string
    Age  int
    WNI  bool
}
```

Cara membuat object dari struct:

```go
p := Person{
    Name: "John",
    Age:  17,
    WNI:  true,
}

fmt.Println(p)
```

Bisa juga tanpa menulis nama field, tetapi urutannya harus sama seperti definisi struct:

```go
p := Person{"John", 17, true}
fmt.Println(p)
```

Catatan:

Default value di Go bukan `null`, tapi zero value.

Contoh zero value:

- `string` = `""`
- `int` = `0`
- `bool` = `false`

Jadi kalau struct dibuat tanpa value, field-nya akan otomatis berisi zero value.

---

## 6. For Loop

Di Go, perulangan hanya menggunakan `for`.

### Infinite Loop

```go
for {
    fmt.Println("Hello")
}
```

Loop di atas akan jalan selamanya dan tidak akan berhenti kecuali dihentikan manual atau pakai `break`.

### For dengan Kondisi

```go
for condition {
    // logic
}
```

Contoh:

```go
i := 0

for i < 10 {
    fmt.Println(i)
    i++
}
```

### For dengan Counter

```go
for i := 0; i < 10; i++ {
    // logic
}
```

### Range untuk Slice

```go
numbers := []int{1, 2, 3, 4, 5}

for i, value := range numbers {
    fmt.Println(i, value)
}
```

`i` adalah index, `value` adalah isi dari slice.

Kalau index tidak dipakai, gunakan blank identifier `_`.

```go
for _, value := range numbers {
    fmt.Println(value)
}
```

---

## 7. If Else

Contoh struktur `if else`:

```go
if x == 3 {
    // logic
} else if x == 2 {
    // logic
} else {
    // logic
}
```

---

## 8. Switch

Switch digunakan untuk memilih logic berdasarkan value tertentu.

```go
switch number {
case 1:
    // logic
case 2:
    // logic
case 3:
    // logic
default:
    // logic
}
```

Di Go, setiap `case` otomatis berhenti setelah logic case itu selesai. Jadi biasanya tidak perlu `break`.

Kalau ingin lanjut ke case berikutnya, bisa pakai `fallthrough`.

```go
switch number {
case 1:
    fmt.Println("satu")
    fallthrough
case 2:
    fmt.Println("dua")
default:
    fmt.Println("default")
}
```

Catatan:

`fallthrough` akan lanjut menjalankan case setelahnya, tapi hanya satu case berikutnya.

---

## 9. Function

Function digunakan untuk membuat logic yang bisa dipakai ulang.

```go
func Greetings(firstName string, lastName string, age int, number int) string {
    return "hello " + firstName + " " + lastName
}
```

Function juga bisa punya lebih dari satu return value.

```go
func Div(a, b int) (int, bool) {
    if b == 0 {
        return 0, false
    }

    return a / b, true
}
```

Contoh pemakaian:

```go
result, ok := Div(10, 2)
fmt.Println(result, ok)
```

Kalau salah satu return value tidak digunakan, bisa pakai `_`.

```go
result, _ := Div(10, 2)
fmt.Println(result)
```

---

## 10. Named Return Value

Go bisa menggunakan named return value, yaitu nama variable return ditulis langsung di signature function.

```go
func Addition(value ...int) (result int) {
    for _, v := range value {
        result += v
    }

    return
}
```

Contoh pemakaian:

```go
numbers := []int{1, 10, 13, 79}

result := Addition(numbers...)
fmt.Println(result)
```

Catatan:

- `value ...int` berarti function menerima jumlah parameter `int` yang fleksibel.
- `numbers...` digunakan untuk mengirim slice sebagai variadic argument.
- Karena return value sudah diberi nama `result`, maka `return` bisa langsung dipakai tanpa menulis `return result`.

---

## 11. Defer

`defer` digunakan untuk menunda eksekusi sampai function selesai.

```go
func main() {
    defer fmt.Println("jalan terakhir")

    fmt.Println("jalan pertama")
}
```

Output:

```text
jalan pertama
jalan terakhir
```

`defer` biasanya dipakai untuk cleanup, misalnya menutup file.

```go
file, err := os.Open("go.mod")
if err != nil {
    fmt.Println("error", err)
    return
}
defer file.Close()
```

Catatan:

Kalau ada lebih dari satu `defer`, Go menjalankannya dengan prinsip stack: yang terakhir ditulis akan dijalankan lebih dulu.

---

## 12. Pointer dan Passing by Reference

Kalau struct dikirim sebagai parameter biasa, Go akan membuat copy dari object tersebut. Jadi perubahan di dalam function tidak mengubah data aslinya.

```go
func increaseAge(p Person) {
    p.Age++
}
```

Supaya perubahan berdampak ke object asli, gunakan pointer.

```go
func increaseAge(p *Person) {
    p.Age++
}
```

Contoh pemakaian:

```go
p := Person{
    Name: "John",
    Age:  18,
}

increaseAge(&p)
fmt.Println(p.Age)
```

Catatan:

- `&p` mengambil alamat memory dari variable `p`.
- `*Person` berarti parameter menerima pointer ke `Person`.
- Untuk akses field struct dari pointer, Go bisa langsung pakai `p.Age`.

---

## 13. Go Module dan Package

`go mod init` digunakan untuk menginisialisasi project Go yang menggunakan Go Modules.

Contoh:

```bash
go mod init github.com/yeka/crayon
```

Kalau module didefinisikan sebagai:

```text
github.com/yeka/crayon
```

Maka package di dalam project bisa di-import seperti ini:

```go
import "github.com/yeka/crayon/repo"
```

Contoh pemakaian struct dari package lain:

```go
p := repo.Person{
    Name: "Jon",
    Age:  18,
}
```

Catatan tentang public/private:

Di Go, nama yang diawali huruf besar bersifat public/exported, sehingga bisa diakses dari package lain.

```go
type Person struct {
    Name string // public
    age  int    // private, hanya bisa diakses dari package yang sama
}
```

Nama function juga sama:

```go
func PublicFunction() {
    // bisa diakses dari package lain
}

func privateFunction() {
    // hanya bisa diakses dari package yang sama
}
```

---

## 14. Function sebagai Type

Function bisa dijadikan tipe data.

```go
type Operator func(a, b int) int
```

Contoh:

```go
func Add(a, b int) int {
    return a + b
}

func Sub(a, b int) int {
    return a - b
}

func Calculate(a, b int, op Operator) int {
    return op(a, b)
}
```

Pemakaian:

```go
func main() {
    var op Operator

    op = Add
    result := Calculate(1, 2, op)

    fmt.Println(result)
}
```

Function yang dimasukkan ke `Calculate` harus punya signature yang sama dengan `Operator`, yaitu:

```go
func(int, int) int
```

---

## 15. Method / Function Receiver

Go tidak punya class seperti OOP biasa, tapi bisa membuat method pada struct menggunakan receiver.

```go
type Person struct {
    Name string
    Age  int
}

func (p *Person) IncreaseAge() {
    p.Age++
}
```

Pemakaian:

```go
person := Person{
    Name: "John",
    Age:  18,
}

person.IncreaseAge()
fmt.Println(person.Age)
```

Catatan:

Method receiver ini mirip seperti method milik object. Jadi daripada memanggil:

```go
IncreaseAge(&person)
```

Kita bisa memanggil:

```go
person.IncreaseAge()
```

---

## 16. Interface

Interface berisi kumpulan method yang harus dimiliki oleh suatu tipe.

```go
type AgeIncreaser interface {
    IncreaseAge()
}
```

Kalau suatu tipe punya method `IncreaseAge()`, maka tipe itu otomatis memenuhi interface `AgeIncreaser`.

```go
func IncreaseAge(v AgeIncreaser) {
    v.IncreaseAge()
}
```

Contoh:

```go
type Person struct {
    Name string
    Age  int
}

func (p *Person) IncreaseAge() {
    p.Age++
}

func main() {
    p := &Person{
        Name: "John",
        Age:  18,
    }

    IncreaseAge(p)
    fmt.Println(p.Age)
}
```

Catatan:

Di Go, tidak perlu menulis `implements`. Kalau method-nya cocok dengan interface, maka otomatis compatible.

Interface kosong juga bisa dibuat:

```go
func PrintAnything(v interface{}) {
    fmt.Println(v)
}
```

Sejak Go versi baru, `interface{}` bisa juga ditulis sebagai `any`.

```go
func PrintAnything(v any) {
    fmt.Println(v)
}
```

---

## 17. File I/O: Baca File

Contoh membaca file menggunakan `os.Open` dan `Read`.

```go
func main() {
    f, err := os.Open("go.mod")
    if err != nil {
        fmt.Println("error", err.Error())
        return
    }
    defer f.Close()

    b := make([]byte, 1024)

    for {
        n, err := f.Read(b)
        if err == io.EOF {
            break
        }

        if err != nil {
            fmt.Println("read error", err)
            return
        }

        fmt.Println("Jumlah byte terbaca:", n)
        fmt.Println("Content:")
        fmt.Println(string(b[:n]))
    }
}
```

Catatan:

- `os.Open()` membuka file.
- `defer f.Close()` memastikan file ditutup setelah function selesai.
- `make([]byte, 1024)` membuat buffer 1024 byte.
- `f.Read(b)` membaca isi file ke buffer.
- `io.EOF` berarti file sudah selesai dibaca.
- `string(b[:n])` mengubah byte yang terbaca menjadi string.

Cara lain yang lebih simpel adalah menggunakan `io.ReadAll`.

```go
func main() {
    f, err := os.Open("go.mod")
    if err != nil {
        fmt.Println("error", err.Error())
        return
    }
    defer f.Close()

    b, err := io.ReadAll(f)
    if err != nil {
        fmt.Println("read error", err)
        return
    }

    fmt.Println("Jumlah byte terbaca:", len(b))
    fmt.Println("Content:")
    fmt.Println(string(b))
}
```

---

## 18. File I/O: Tulis File

Contoh membuat dan menulis file:

```go
func WriteFileExample() {
    f, err := os.OpenFile(
        "b.txt",
        os.O_CREATE|os.O_WRONLY,
        0655,
    )

    if err != nil {
        fmt.Println("error", err)
        return
    }
    defer f.Close()

    n, err := f.Write([]byte("Hello World!"))
    if err != nil {
        fmt.Println("write error", err)
        return
    }

    fmt.Println("Byte written:", n)
}
```

Catatan:

- `os.O_CREATE` membuat file kalau belum ada.
- `os.O_WRONLY` membuka file untuk ditulis.
- `0655` adalah permission file.
- `f.Write()` menerima data dalam bentuk `[]byte`.

---

## 19. fmt.Fprintf dengan io.Writer

`fmt.Fprintf` bisa menerima parameter berupa `io.Writer`.

Contoh:

```go
fmt.Fprintf(f, "hello world")
```

Karena file memenuhi interface `io.Writer`, maka file bisa menjadi target output untuk `fmt.Fprintf`.

---

## 20. Embed

`embed` digunakan untuk memasukkan file ke dalam binary Go saat compile time.

Contoh:

```go
package main

import (
    "embed"
    "fmt"
)

//go:embed a.txt
var content string

func main() {
    fmt.Println(content)
}
```

Catatan:

- Package `embed` harus di-import.
- Directive `//go:embed a.txt` digunakan untuk memasukkan isi file `a.txt`.
- Variable di bawah directive akan berisi isi file tersebut.

Kalau package `embed` tidak dipakai langsung di code, bisa import dengan blank identifier.

```go
import _ "embed"
```

---

## 21. Concurrency dan Goroutine

Goroutine digunakan untuk menjalankan function secara concurrent.

Tanpa goroutine:

```go
func main() {
    start := time.Now()
    defer func() {
        fmt.Println(time.Since(start))
    }()

    for i := range 4 {
        SlowProcess(i)
    }
}

func SlowProcess(i int) {
    fmt.Println(i)
    time.Sleep(2 * time.Second)
}
```

Kode di atas menjalankan process satu per satu, jadi total waktunya lebih lama.

Dengan goroutine:

```go
func main() {
    start := time.Now()
    defer func() {
        fmt.Println(time.Since(start))
    }()

    for i := range 4 {
        go SlowProcess(i)
    }
}

func SlowProcess(i int) {
    fmt.Println(i)
    time.Sleep(2 * time.Second)
}
```

Goroutine membuat `SlowProcess(i)` dijalankan di tempat berbeda secara concurrent.

Namun kalau hanya memakai `go SlowProcess(i)`, program utama bisa selesai lebih dulu sebelum goroutine selesai. Karena itu perlu mekanisme tunggu seperti `sync.WaitGroup`.

---

## 22. WaitGroup

`sync.WaitGroup` digunakan untuk menunggu semua goroutine selesai.

```go
func main() {
    start := time.Now()
    defer func() {
        fmt.Println(time.Since(start))
    }()

    wg := sync.WaitGroup{}

    for i := range 4 {
        wg.Add(1)
        go SlowProcess(i, &wg)
    }

    fmt.Println("waiting")
    wg.Wait()
}

func SlowProcess(i int, wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Println(i)
    time.Sleep(2 * time.Second)
}
```

Penjelasan:

- `wg := sync.WaitGroup{}` membuat WaitGroup.
- `wg.Add(1)` menambah jumlah goroutine yang harus ditunggu.
- `wg.Done()` menandakan satu goroutine sudah selesai.
- `wg.Wait()` menunggu sampai semua goroutine selesai.
- `defer wg.Done()` aman dipakai supaya `Done()` tetap terpanggil saat function selesai.

---

## 23. Channel

Channel digunakan untuk komunikasi antar goroutine.

Membuat channel:

```go
ch := make(chan int)
```

Mengirim data ke channel:

```go
ch <- 1
```

Menerima data dari channel:

```go
val := <-ch
```

Contoh sederhana:

```go
func Channel() {
    ch := make(chan int)

    go func() {
        time.Sleep(3 * time.Second)

        val := <-ch
        fmt.Println("read from channel", val)
    }()

    fmt.Println("waiting until something read from channel")

    ch <- 1

    fmt.Println("done")
}
```

Catatan:

- Channel biasa bersifat blocking.
- `ch <- 1` akan menunggu sampai ada goroutine lain yang menerima data.
- `<-ch` akan menunggu sampai ada data yang dikirim.
- Channel sering dipakai untuk sinkronisasi dan komunikasi antar goroutine.

---

## 24. Ringkasan Konsep Penting

### Zero Value

Go punya default value untuk setiap tipe data:

```go
var name string // ""
var age int     // 0
var active bool // false
```

### Blank Identifier

Gunakan `_` kalau ada value yang tidak mau dipakai.

```go
for _, value := range numbers {
    fmt.Println(value)
}
```

Atau saat function punya beberapa return value:

```go
result, _ := Div(10, 2)
```

### Public dan Private

Nama yang diawali huruf besar bisa diakses dari package lain.

```go
func PublicFunction() {}
```

Nama yang diawali huruf kecil hanya bisa diakses dari package yang sama.

```go
func privateFunction() {}
```

### Pointer

Gunakan pointer kalau function perlu mengubah data asli.

```go
func Update(p *Person) {
    p.Age++
}
```

### Interface

Interface di Go implicit. Tidak perlu keyword `implements`.

```go
type Swimmer interface {
    Swim()
}
```

Kalau suatu tipe punya method `Swim()`, maka otomatis dianggap memenuhi interface `Swimmer`.

### Goroutine

Gunakan `go` untuk menjalankan function secara concurrent.

```go
go doSomething()
```

### WaitGroup

Gunakan `WaitGroup` kalau program perlu menunggu semua goroutine selesai.

```go
wg.Add(1)
go func() {
    defer wg.Done()
    // logic
}()
wg.Wait()
```

### Channel

Gunakan channel untuk mengirim dan menerima data antar goroutine.

```go
ch := make(chan string)

go func() {
    ch <- "hello"
}()

msg := <-ch
fmt.Println(msg)
```
