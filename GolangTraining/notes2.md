# Channel

Channel digunakan untuk komunikasi antar goroutine.

Channel bisa di-iterate menggunakan `for range`.

```go
for v := range ch {
    fmt.Println(v)
}
```

Artinya: selama channel `ch` masih terbuka dan masih ada data yang dikirim, data akan diterima satu per satu dan dimasukkan ke variable `v`.

Contoh:

```go
func Channel() {
    ch := make(chan int)

    go func() {
        for v := range ch {
            fmt.Println("read from channel", v)
        }

        fmt.Println("channel closed")
    }()

    fmt.Println("waiting until something read from channel")

    ch <- 1
    close(ch)

    fmt.Println("done")
}
```

Penjelasan:

- `ch := make(chan int)` membuat channel bertipe `int`.
- `ch <- 1` mengirim data ke channel.
- `for v := range ch` menerima data dari channel selama channel masih terbuka.
- `close(ch)` menandakan bahwa tidak akan ada data lagi yang dikirim ke channel.
- Setelah channel di-close, perulangan `for range` akan berhenti.

Catatan penting:

Channel perlu di-close kalau receiver menggunakan `for range`, karena `for range` akan terus menunggu data baru selama channel belum ditutup.

Bukan karena channel harus selalu di-close supaya memory hilang, tapi karena `close(ch)` memberi sinyal ke receiver bahwa pengiriman data sudah selesai.

Kalau channel tidak di-close, bagian ini bisa menunggu terus:

```go
for v := range ch {
    fmt.Println(v)
}
```

Yang biasanya melakukan `close(ch)` adalah sender, bukan receiver. Receiver cukup membaca data dari channel.

---

## Channel dengan WaitGroup

Kalau menggunakan goroutine, program utama bisa selesai duluan sebelum goroutine selesai menjalankan logic-nya.

Karena itu, kita bisa menggunakan `sync.WaitGroup` untuk menunggu goroutine selesai.

```go
func Channel() {
    ch := make(chan int)
    wg := sync.WaitGroup{}

    wg.Add(1)

    go func() {
        defer wg.Done()

        for v := range ch {
            fmt.Println("read from channel", v)
        }

        fmt.Println("channel closed")
    }()

    fmt.Println("waiting until something read from channel")

    ch <- 1
    close(ch)

    wg.Wait()

    fmt.Println("done")
}
```

Penjelasan:

- `wg := sync.WaitGroup{}` membuat WaitGroup.
- `wg.Add(1)` menandakan ada 1 goroutine yang perlu ditunggu.
- `defer wg.Done()` menandakan goroutine sudah selesai saat function berakhir.
- `wg.Wait()` membuat main goroutine menunggu sampai semua goroutine selesai.
- `close(ch)` membuat `for range ch` berhenti.
- Setelah goroutine selesai, baru program lanjut ke `fmt.Println("done")`.

Kenapa perlu `WaitGroup`?

Karena goroutine jalan secara concurrent. Jadi tanpa `wg.Wait()`, bisa saja `main` function selesai duluan sebelum goroutine sempat menyelesaikan proses membaca channel.

Menggunakan `WaitGroup` adalah best practice kalau kita perlu memastikan semua goroutine selesai sebelum program lanjut atau selesai.

---

## Contoh dengan Banyak Data

```go
func Channel() {
    ch := make(chan int)
    wg := sync.WaitGroup{}

    wg.Add(1)

    go func() {
        defer wg.Done()

        for v := range ch {
            fmt.Println("read from channel", v)
        }

        fmt.Println("channel closed")
    }()

    for i := 1; i <= 5; i++ {
        ch <- i
    }

    close(ch)

    wg.Wait()

    fmt.Println("done")
}
```

Output kira-kira:

```text
read from channel 1
read from channel 2
read from channel 3
read from channel 4
read from channel 5
channel closed
done
```

Intinya:

- `for range ch` dipakai untuk menerima data channel berkali-kali.
- `close(ch)` dipakai untuk menghentikan loop `for range`.
- `WaitGroup` dipakai supaya program utama menunggu goroutine selesai.

---

# Concurrent Worker

Concurrent worker digunakan untuk membatasi jumlah goroutine yang bekerja secara bersamaan.

Kalau kita membuat terlalu banyak goroutine sekaligus, program belum tentu jadi lebih cepat. Bisa saja processor malah terlalu sibuk melakukan context switching, yaitu berpindah-pindah dari satu goroutine ke goroutine lain.

Karena itu, untuk pekerjaan yang banyak, biasanya kita buat jumlah worker yang terbatas.

Konsepnya:

1. Buat channel untuk menampung job.
2. Jalankan beberapa worker.
3. Worker membaca job dari channel.
4. Kirim semua job ke channel.
5. Setelah semua job dikirim, close channel.
6. Tunggu semua worker selesai menggunakan `WaitGroup`.

---

## Concurrent Worker dengan WaitGroup Add dan Done

```go
func ConcurrentWorker() {
    jobCh := make(chan int)
    wg := sync.WaitGroup{}

    workerCount := 4
    jobCount := 15

    // Worker preparation
    for i := 0; i < workerCount; i++ {
        workerID := i

        wg.Add(1)

        go func() {
            defer wg.Done()

            for job := range jobCh {
                log.Printf("job %v is done by worker %v\n", job, workerID)
            }
 
            log.Printf("worker %v is done\n", workerID)
        }()
    }

    // Feed pekerjaan ke channel
    for i := 0; i < jobCount; i++ {
        jobCh <- i
    }

    // Setelah semua job dikirim, close channel
    close(jobCh)

    // Tunggu semua worker selesai
    wg.Wait()

    log.Println("all jobs are done")
}
```

Penjelasan:

- `jobCh := make(chan int)` membuat channel untuk mengirim job.
- `workerCount := 4` berarti hanya ada 4 worker yang bekerja.
- `jobCount := 15` berarti ada 15 pekerjaan yang akan diproses.
- Worker dibuat lebih dulu menggunakan loop.
- Setiap worker menjalankan goroutine.
- Setiap worker membaca job dari `jobCh` menggunakan `for job := range jobCh`.
- Kalau ada job masuk, salah satu worker akan mengambil dan memproses job tersebut.
- Setelah semua job dikirim, `close(jobCh)` dipanggil.
- Ketika `jobCh` di-close, semua worker akan keluar dari loop `for range`.
- `wg.Wait()` memastikan program menunggu sampai semua worker selesai.

Catatan:

`jobCh <- i` tidak harus selalu `int`. Data yang dikirim ke channel bisa berupa struct juga.

Contoh:

```go
type Person struct {
    Name string
    Age  int
}
```

Lalu channel-nya bisa dibuat seperti ini:

```go
jobCh := make(chan Person)
```

---

## Kenapa Concurrent Worker Berguna?

Misalnya ada 15 job.

Kalau kita membuat 15 goroutine langsung, semua job bisa jalan bersamaan, tapi kalau jumlah job sangat banyak, misalnya 10.000 job, itu bisa terlalu berat.

Dengan worker, kita bisa membatasi jumlah pekerjaan yang berjalan bersamaan.

Contoh:

```go
workerCount := 4
jobCount := 15
```

Artinya:

- Total pekerjaan ada 15.
- Yang bekerja bersamaan maksimal hanya 4 worker.
- Kalau ada job yang selesai, worker itu akan mengambil job berikutnya.

Jadi worker bekerja seperti antrian.

---

## Concurrent Worker dengan WaitGroup.Go

Pada Go versi baru, `sync.WaitGroup` punya method `Go`.

Dengan `wg.Go`, kita tidak perlu menulis manual:

```go
wg.Add(1)

go func() {
    defer wg.Done()

    // logic
}()
```

Karena sudah diringkas menjadi:

```go
wg.Go(func() {
    // logic
})
```

Contoh:

```go
func ConcurrentWorker() {
    jobCh := make(chan int)
    wg := sync.WaitGroup{}

    workerCount := 4
    jobCount := 15

    // Worker preparation
    for i := 0; i < workerCount; i++ {
        workerID := i

        wg.Go(func() {
            for job := range jobCh {
                log.Printf("job %v is done by worker %v\n", job, workerID)
            }

            log.Printf("worker %v is done\n", workerID)
        })
    }

    // Feed pekerjaan ke channel
    for i := 0; i < jobCount; i++ {
        jobCh <- i
    }

    // Setelah semua job dikirim, close channel
    close(jobCh)

    // Tunggu semua worker selesai
    wg.Wait()

    log.Println("all jobs are done")
}
```

Catatan:

- Namanya `wg.Go`, bukan `wg.GO`.
- `wg.Go` otomatis menjalankan function dalam goroutine.
- `wg.Go` otomatis menambahkan task ke WaitGroup.
- Ketika function selesai, task otomatis dianggap selesai.
- Kalau Go version belum support `wg.Go`, gunakan versi `wg.Add(1)` dan `defer wg.Done()`.

---

## Inti Concurrent Worker

Concurrent worker dipakai untuk membatasi jumlah goroutine yang memproses pekerjaan.

Pattern umumnya:

```go
jobCh := make(chan Job)
wg := sync.WaitGroup{}

// Start workers
for i := 0; i < workerCount; i++ {
    wg.Add(1)

    go func() {
        defer wg.Done()

        for job := range jobCh {
            // process job
        }
    }()
}

// Send jobs
for _, job := range jobs {
    jobCh <- job
}

// Stop workers
close(jobCh)

// Wait workers
wg.Wait()
```

Intinya:

- Worker dibuat dulu.
- Job dikirim lewat channel.
- Worker mengambil job dari channel.
- `close(jobCh)` memberi tahu worker bahwa job sudah habis.
- `wg.Wait()` menunggu semua worker selesai.

## Worker Function untuk Concurrent Worker

Worker tidak harus ditulis langsung di dalam `ConcurrentWorker`. Kalau mau lebih rapi, logic worker bisa dipisah ke function sendiri.

Contoh function worker:

```go
func Worker(workerID int, wg *sync.WaitGroup, jobCh <-chan int) {
    defer wg.Done()

    for job := range jobCh {
        log.Printf("job %v is done by worker %v\n", job, workerID)
    }

    log.Printf("worker %v is done\n", workerID)
}
```

Penjelasan parameter:

```go
workerID int
```

Digunakan untuk memberi identitas ke worker. Misalnya worker `0`, worker `1`, worker `2`, dan seterusnya.

```go
wg *sync.WaitGroup
```

Digunakan supaya worker bisa memberi tahu bahwa pekerjaannya sudah selesai.

Karena `WaitGroup` perlu dimodifikasi dari dalam function worker, maka dikirim sebagai pointer menggunakan `*sync.WaitGroup`.

```go
jobCh <-chan int
```

Artinya `jobCh` adalah receive-only channel.

Receive-only channel berarti di dalam function `Worker`, channel ini hanya boleh digunakan untuk menerima data, bukan mengirim data.

Contoh yang boleh:

```go
job := <-jobCh
```

Atau:

```go
for job := range jobCh {
    // process job
}
```

Contoh yang tidak boleh:

```go
jobCh <- 1
```

Karena `jobCh <-chan int` hanya bisa menerima data.

---

## Arah Channel di Go

Ada 3 bentuk penulisan channel:

```go
jobCh chan int
```

Artinya channel bisa digunakan untuk mengirim dan menerima data.

```go
jobCh <-chan int
```

Artinya channel hanya bisa digunakan untuk menerima data.

```go
jobCh chan<- int
```

Artinya channel hanya bisa digunakan untuk mengirim data.

Jadi penulisan yang benar untuk receive-only channel adalah:

```go
jobCh <-chan int
```

Bukan:

```go
jobCh <- chan int
```

Dan bukan juga:

```go
<- jobCh chan int
```

---

## Concurrent Worker dengan Function Worker

Contoh lengkap:

```go
func Worker(workerID int, wg *sync.WaitGroup, jobCh <-chan int) {
    defer wg.Done()

    for job := range jobCh {
        log.Printf("job %v is done by worker %v\n", job, workerID)
    }

    log.Printf("worker %v is done\n", workerID)
}

func ConcurrentWorker() {
    jobCh := make(chan int)
    wg := sync.WaitGroup{}

    workerCount := 4
    jobCount := 15

    // Start workers
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go Worker(i, &wg, jobCh)
    }

    // Send jobs
    for i := 0; i < jobCount; i++ {
        jobCh <- i
    }

    // Stop workers
    close(jobCh)

    // Wait until all workers are done
    wg.Wait()

    log.Println("all jobs are done")
}
```

Penjelasan alur:

1. `jobCh := make(chan int)` membuat channel untuk mengirim job.
2. `wg := sync.WaitGroup{}` membuat WaitGroup untuk menunggu worker selesai.
3. Loop pertama membuat beberapa worker.
4. Setiap worker dijalankan sebagai goroutine.
5. Setiap worker membaca job dari `jobCh`.
6. Loop kedua mengirim semua job ke `jobCh`.
7. Setelah semua job dikirim, `close(jobCh)` dipanggil.
8. Ketika `jobCh` di-close, semua worker akan keluar dari `for job := range jobCh`.
9. `wg.Wait()` menunggu semua worker selesai.
10. Setelah semua worker selesai, program lanjut ke `log.Println("all jobs are done")`.

---

## Kenapa `jobCh` Dibuat Receive-Only di Worker?

Di function `Worker`, tugas worker hanya mengambil job dari channel, bukan mengirim job baru ke channel.

Karena itu, parameter ini lebih aman ditulis sebagai receive-only channel:

```go
jobCh <-chan int
```

Dengan begitu, kalau di dalam `Worker` kita tidak sengaja menulis:

```go
jobCh <- 1
```

Program akan error saat compile.

Ini bagus karena Go membantu mencegah kesalahan penggunaan channel.

---

## Versi dengan `WaitGroup.Go`

Kalau Go version yang digunakan sudah mendukung `wg.Go`, function `Worker` tidak perlu menerima `WaitGroup`.

Function worker cukup menerima `workerID` dan `jobCh`.

```go
func Worker(workerID int, jobCh <-chan int) {
    for job := range jobCh {
        log.Printf("job %v is done by worker %v\n", job, workerID)
    }

    log.Printf("worker %v is done\n", workerID)
}
```

Contoh lengkap:

```go
func Worker(workerID int, jobCh <-chan int) {
    for job := range jobCh {
        log.Printf("job %v is done by worker %v\n", job, workerID)
    }

    log.Printf("worker %v is done\n", workerID)
}

func ConcurrentWorker() {
    jobCh := make(chan int)
    wg := sync.WaitGroup{}

    workerCount := 4
    jobCount := 15

    // Start workers
    for i := 0; i < workerCount; i++ {
        workerID := i

        wg.Go(func() {
            Worker(workerID, jobCh)
        })
    }

    // Send jobs
    for i := 0; i < jobCount; i++ {
        jobCh <- i
    }

    // Stop workers
    close(jobCh)

    // Wait until all workers are done
    wg.Wait()

    log.Println("all jobs are done")
}
```

Dengan `wg.Go`, kita tidak perlu menulis manual:

```go
wg.Add(1)

go func() {
    defer wg.Done()

    // logic
}()
```

Karena `wg.Go` sudah menangani proses menambahkan task ke WaitGroup, menjalankan goroutine, dan menandai task selesai ketika function selesai.

---

## Intinya

Function worker bisa ditulis seperti ini:

```go
func Worker(workerID int, wg *sync.WaitGroup, jobCh <-chan int)
```

Kalau menggunakan `wg.Go`, bisa ditulis lebih sederhana:

```go
func Worker(workerID int, jobCh <-chan int)
```

Gunakan:

```go
jobCh <-chan int
```

kalau function hanya perlu menerima data dari channel.

Gunakan:

```go
jobCh chan<- int
```

kalau function hanya perlu mengirim data ke channel.

Gunakan:

```go
jobCh chan int
```

kalau function perlu mengirim dan menerima data dari channel. 

Kalau ada 100 worker dan 100000 job, bukan berarti 100000 job jalan bersamaan.

Yang jalan bersamaan maksimal hanya 100 job, karena jumlah worker-nya 100.

Kalau semua worker sedang sibuk, pengiriman job berikutnya ke channel akan blocking sampai ada worker yang selesai dan siap menerima job baru.

Dengan unbuffered channel:

```go
jobCh := make(chan int)
```

# Mutex

`Mutex` digunakan untuk mengamankan data yang diakses oleh banyak goroutine secara bersamaan.

Mutex bisa dipakai di concurrent worker maupun concurrent biasa.

Contoh kasus yang sering butuh mutex adalah ketika banyak goroutine mengubah variable yang sama.

Misalnya ada slice:

```go
func Mutex() {
    result := []int{}

    for i := 0; i < 10; i++ {
        result = append(result, i)
    }

    fmt.Println(len(result))
}
```

Kalau kode di atas dijalankan biasa tanpa goroutine, hasilnya aman karena prosesnya berjalan satu per satu.

Tapi kalau pakai goroutine:

```go
func Mutex() {
    result := []int{}

    for i := 0; i < 100; i++ {
        go func() {
            result = append(result, i)
        }()
    }

    fmt.Println(len(result))
}
```

Kode di atas bermasalah karena banyak goroutine mengakses dan mengubah `result` secara bersamaan.

Hasilnya bisa tidak sesuai. Misalnya harusnya `len(result)` adalah `100`, tapi yang keluar bisa `82`, `91`, atau angka lain.

Masalah seperti ini disebut **race condition** atau **data race**.

---

## WaitGroup Tidak Menyelesaikan Race Condition

Kita bisa menambahkan `WaitGroup` supaya program menunggu semua goroutine selesai.

```go
func Mutex() {
    result := []int{}
    wg := sync.WaitGroup{}

    for i := 0; i < 100; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()

            result = append(result, i)
        }()
    }

    wg.Wait()

    fmt.Println(len(result))
}
```

Tapi kode di atas masih belum aman.

`WaitGroup` hanya memastikan semua goroutine selesai sebelum program lanjut.

`WaitGroup` tidak mencegah beberapa goroutine mengubah `result` secara bersamaan.

Jadi walaupun sudah memakai `WaitGroup`, hasilnya tetap bisa salah karena masih ada race condition.

---

## Menggunakan Mutex

Untuk mencegah race condition, kita bisa memakai `sync.Mutex`.

```go
mutex := sync.Mutex{}
```

Mutex punya dua method utama:

```go
mutex.Lock()
mutex.Unlock()
```

Artinya:

- `mutex.Lock()` mengunci akses ke bagian tertentu.
- `mutex.Unlock()` membuka kunci setelah bagian tersebut selesai.
- Selama satu goroutine sedang memegang lock, goroutine lain harus menunggu.

Contoh:

```go
func Mutex() {
    result := []int{}
    wg := sync.WaitGroup{}
    mutex := sync.Mutex{}

    for i := 0; i < 100; i++ {
        i := i

        wg.Add(1)

        go func() {
            defer wg.Done()

            mutex.Lock()
            result = append(result, i)
            mutex.Unlock()
        }()
    }

    wg.Wait()

    fmt.Println(len(result))
}
```

Dengan mutex, hanya satu goroutine yang boleh melakukan `append` ke `result` dalam satu waktu.

Jadi `result = append(result, i)` menjadi aman.

---

## Menggunakan defer untuk Unlock

Biasanya `mutex.Unlock()` ditulis dengan `defer`, supaya lock tetap terbuka walaupun function selesai karena error atau return lebih awal.

```go
func Mutex() {
    result := []int{}
    wg := sync.WaitGroup{}
    mutex := sync.Mutex{}

    for i := 0; i < 100; i++ {
        i := i

        wg.Add(1)

        go func() {
            defer wg.Done()

            mutex.Lock()
            defer mutex.Unlock()

            result = append(result, i)
        }()
    }

    wg.Wait()

    fmt.Println(len(result))
}
```

Penjelasan:

```go
mutex.Lock()
defer mutex.Unlock()
```

Artinya goroutine akan mengunci dulu, lalu `mutex.Unlock()` pasti dipanggil saat goroutine selesai.

---

## Lock Hanya Bagian yang Perlu Diamankan

Mutex memang membuat data aman, tapi jangan terlalu banyak logic dimasukkan di antara `Lock()` dan `Unlock()`.

Contoh yang kurang bagus:

```go
go func() {
    defer wg.Done()

    mutex.Lock()
    defer mutex.Unlock()

    time.Sleep(2 * time.Second)

    result = append(result, i)
}()
```

Kode di atas kurang bagus karena `time.Sleep` ada di dalam lock.

Akibatnya, selama proses lambat itu berjalan, goroutine lain tidak bisa mengakses `result`.

Kalau bagian yang lambat dimasukkan ke dalam lock, concurrency jadi percuma, karena semua goroutine harus menunggu satu per satu terlalu lama.

Yang lebih baik:

```go
go func() {
    defer wg.Done()

    time.Sleep(2 * time.Second)

    mutex.Lock()
    result = append(result, i)
    mutex.Unlock()
}()
```

Di sini proses lambat dilakukan di luar lock.

Lock hanya dipakai saat benar-benar perlu mengubah shared data, yaitu saat `append` ke `result`.

---

## Contoh Final

```go
func Mutex() {
    result := []int{}
    wg := sync.WaitGroup{}
    mutex := sync.Mutex{}

    for i := 0; i < 100; i++ {
        i := i

        wg.Add(1)

        go func() {
            defer wg.Done()

            // Simulasi proses yang bisa berjalan concurrent
            time.Sleep(100 * time.Millisecond)

            // Bagian ini perlu diamankan karena mengubah shared data
            mutex.Lock()
            result = append(result, i)
            mutex.Unlock()
        }()
    }

    wg.Wait()

    fmt.Println(len(result))
}
```

Penjelasan:

- `result` adalah shared data karena diakses banyak goroutine.
- `wg` digunakan untuk menunggu semua goroutine selesai.
- `mutex` digunakan supaya hanya satu goroutine yang bisa mengubah `result` dalam satu waktu.
- `time.Sleep` diletakkan di luar lock supaya proses concurrent tetap efektif.
- `mutex.Lock()` dan `mutex.Unlock()` hanya membungkus bagian yang benar-benar mengubah data.

---

## Intinya

`WaitGroup` digunakan untuk menunggu goroutine selesai.

`Mutex` digunakan untuk mencegah race condition ketika banyak goroutine mengakses atau mengubah data yang sama.

Kalau hanya pakai `WaitGroup`, program memang menunggu semua goroutine selesai, tapi data tetap bisa rusak karena diubah bersamaan.

Kalau pakai `Mutex`, akses ke shared data jadi aman karena hanya satu goroutine yang boleh masuk ke critical section dalam satu waktu.

Critical section adalah bagian kode yang mengakses atau mengubah shared data.

Contoh critical section:

```go
mutex.Lock()
result = append(result, i)
mutex.Unlock()
```

Best practice:

- Pakai `WaitGroup` untuk menunggu goroutine selesai.
- Pakai `Mutex` untuk melindungi shared data.
- Jangan taruh proses lambat di dalam `Lock()` dan `Unlock()`.
- Lock hanya bagian yang benar-benar perlu diamankan.

# Time

Package `time` digunakan untuk bekerja dengan waktu, tanggal, durasi, dan parsing string menjadi tipe `time.Time`.

Di Go, format untuk parsing waktu agak unik. Go tidak memakai format seperti `YYYY-MM-DD`, tapi memakai reference time khusus:

```go
2006-01-02 15:04:05
```

Reference time ini harus diingat karena Go memakai angka tersebut sebagai patokan format.

---

## time.Parse

`time.Parse` digunakan untuk mengubah string menjadi `time.Time`.

Contoh:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    layout := "2006-01-02 15:04:05"
    value := "2026-06-23 14:30:00"

    parsedTime, err := time.Parse(layout, value)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    fmt.Println(parsedTime)
}
```

Penjelasan:

```go
layout := "2006-01-02 15:04:05"
```

Artinya format waktu yang diharapkan adalah:

```text
tahun-bulan-tanggal jam:menit:detik
```

Contoh value yang cocok:

```text
2026-06-23 14:30:00
```

---

## Format 24 Jam

Kalau ingin memakai format 24 jam, gunakan `15`.

```go
layout := "2006-01-02 15:04:05"
value := "2026-06-23 21:15:30"
```

Di sini `21` berarti jam 9 malam.

---

## Format 12 Jam dengan AM/PM

Kalau ingin memakai format 12 jam, gunakan `03` dan tambahkan `PM`.

```go
layout := "2006-01-02 03:04:05PM"
value := "2026-06-23 09:15:30PM"
```

Kalau string waktunya punya spasi sebelum AM/PM, layout-nya juga harus punya spasi.

```go
layout := "2006-01-02 03:04:05 PM"
value := "2026-06-23 09:15:30 PM"
```

Catatan:

- `15` digunakan untuk format 24 jam.
- `03` digunakan untuk format 12 jam.
- `PM` digunakan untuk membaca AM/PM.
- Yang benar adalah `2006`, bukan `20060`.

---

# Cobra

`cobra` adalah library Go yang digunakan untuk membuat CLI atau command-line application.

CLI adalah aplikasi yang dijalankan melalui terminal.

Contoh command-line application:

```bash
app
app server
app worker
app migrate
```

Dengan `cobra`, kita bisa membuat:

- command utama
- subcommand
- argument
- flag
- logic yang dijalankan saat command dipanggil

---

## Install Cobra

Untuk menggunakan Cobra, jalankan:

```bash
go get github.com/spf13/cobra
```

---

## Command Dasar Cobra

Contoh command sederhana:

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
)

func main() {
    cmd := cobra.Command{
        Use: "app",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello world!")
        },
    }

    cmd.Execute()
}
```

Penjelasan:

```go
cmd := cobra.Command{
    Use: "app",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}
```

Artinya kita membuat command utama bernama `app`.

Bagian ini:

```go
Use: "app"
```

digunakan untuk menentukan nama command.

Bagian ini:

```go
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello world!")
}
```

adalah logic yang akan dijalankan saat command dipanggil.

Bagian ini:

```go
cmd.Execute()
```

digunakan untuk menjalankan command.

---

## Error Handling saat Execute

Biasanya `cmd.Execute()` lebih aman ditulis dengan error handling.

```go
func main() {
    cmd := cobra.Command{
        Use: "app",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello world!")
        },
    }

    if err := cmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

---

## Cobra dengan Subcommand

Subcommand adalah command tambahan di bawah command utama.

Misalnya kita punya command utama:

```bash
app
```

Lalu kita ingin punya subcommand:

```bash
app server
app worker
```

Maka kita bisa menggunakan `AddCommand`.

Contoh lengkap:

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
)

func main() {
    rootCmd := cobra.Command{
        Use:   "app",
        Short: "Main application command",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello world!")
        },
    }

    serverCmd := cobra.Command{
        Use:   "server",
        Short: "Running a server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Server is running!")
        },
    }

    workerCmd := cobra.Command{
        Use:   "worker",
        Short: "Running a worker",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Worker is running!")
        },
    }

    rootCmd.AddCommand(&serverCmd, &workerCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Penjelasan:

```go
rootCmd := cobra.Command{
    Use: "app",
}
```

Ini adalah command utama.

```go
serverCmd := cobra.Command{
    Use: "server",
}
```

Ini adalah subcommand untuk menjalankan server.

```go
workerCmd := cobra.Command{
    Use: "worker",
}
```

Ini adalah subcommand untuk menjalankan worker.

```go
rootCmd.AddCommand(&serverCmd, &workerCmd)
```

Artinya `serverCmd` dan `workerCmd` dimasukkan sebagai subcommand dari `rootCmd`.

Jadi command yang bisa dijalankan adalah:

```bash
go run main.go
go run main.go server
go run main.go worker
```

Kalau sudah di-build menjadi binary `app`, command-nya menjadi:

```bash
./app
./app server
./app worker
```

---

## Bentuk AddCommand Langsung

Subcommand juga bisa dibuat langsung di dalam `AddCommand`.

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
)

func main() {
    rootCmd := cobra.Command{
        Use:   "app",
        Short: "Main application command",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello world!")
        },
    }

    rootCmd.AddCommand(
        &cobra.Command{
            Use:   "server",
            Short: "Running a server",
            Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Server is running!")
            },
        },
        &cobra.Command{
            Use:   "worker",
            Short: "Running a worker",
            Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Worker is running!")
            },
        },
    )

    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Ini sama seperti contoh sebelumnya, hanya saja command-nya langsung dibuat di dalam `AddCommand`.

---

## Nested Subcommand

Subcommand juga bisa punya subcommand lagi.

Misalnya:

```bash
app server start
app server stop
```

Berarti `start` dan `stop` adalah subcommand dari `server`.

Contoh:

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
)

func main() {
    rootCmd := cobra.Command{
        Use:   "app",
        Short: "Main application command",
    }

    serverCmd := cobra.Command{
        Use:   "server",
        Short: "Server command",
    }

    startServerCmd := cobra.Command{
        Use:   "start",
        Short: "Start the server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Server is starting!")
        },
    }

    stopServerCmd := cobra.Command{
        Use:   "stop",
        Short: "Stop the server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Server is stopping!")
        },
    }

    serverCmd.AddCommand(&startServerCmd, &stopServerCmd)
    rootCmd.AddCommand(&serverCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Command yang bisa dijalankan:

```bash
go run main.go server start
go run main.go server stop
```

---

## Intinya

`time.Parse` digunakan untuk mengubah string menjadi `time.Time`.

Format waktu di Go menggunakan reference time:

```go
2006-01-02 15:04:05
```

Untuk format 24 jam:

```go
15:04:05
```

Untuk format 12 jam dengan AM/PM:

```go
03:04:05PM
```

`cobra` digunakan untuk membuat CLI.

Pattern dasar Cobra:

```go
cmd := cobra.Command{
    Use: "app",
    Run: func(cmd *cobra.Command, args []string) {
        // logic command
    },
}

cmd.Execute()
```

Untuk menambahkan subcommand, gunakan:

```go
rootCmd.AddCommand(&serverCmd, &workerCmd)
```

Atau langsung:

```go
rootCmd.AddCommand(
    &cobra.Command{
        Use: "server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Server is running!")
        },
    },
)
```
## Struktur Cobra yang Lebih Rapi

Kalau aplikasi CLI sudah mulai punya banyak command dan subcommand, best practice-nya adalah memisahkan command ke file sendiri-sendiri.

Biasanya struktur foldernya seperti ini:

```text
project/
├── main.go
└── cmd/
    ├── root.go
    ├── server.go
    └── worker.go
```

Penjelasan singkat:

- `main.go` hanya memanggil `cmd.Execute()`.
- `cmd/root.go` berisi root command dan function `Execute()`.
- `cmd/server.go` berisi subcommand `server`.
- `cmd/worker.go` berisi subcommand `worker`.
- Setiap subcommand didaftarkan ke root command lewat `func init()`.

---

## main.go

```go
package main

import "app/cmd"

func main() {
    cmd.Execute()
}
```

Penjelasan:

`main.go` dibuat simpel. Dia hanya menjalankan command utama lewat:

```go
cmd.Execute()
```

Jadi logic Cobra tidak ditaruh semua di `main.go`.

---

## cmd/root.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Penjelasan:

```go
var rootCmd = &cobra.Command{
    Use: "app",
}
```

Ini adalah command utama.

Function ini:

```go
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

dipakai supaya dari `main.go`, kita cukup memanggil:

```go
cmd.Execute()
```

Jadi `main.go` tidak perlu tahu detail isi Cobra command-nya.

---

## func init()

Di Go ada function khusus bernama `init`.

```go
func init() {
    // logic initialization
}
```

`init()` akan otomatis dijalankan oleh Go ketika package tersebut di-load atau di-import.

Catatan penting:

- `init()` tidak perlu dipanggil manual.
- `init()` otomatis dijalankan sebelum `main()`.
- `init()` biasanya dipakai untuk setup awal.
- Dalam Cobra, `init()` sering dipakai untuk mendaftarkan subcommand.
- `init()` hanya dijalankan sekali dalam satu aplikasi berjalan.

Contoh:

```go
func init() {
    rootCmd.AddCommand(serverCmd)
}
```

Artinya, saat package `cmd` di-load, subcommand `serverCmd` otomatis didaftarkan ke `rootCmd`.

---

## cmd/server.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running!")
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
}
```

Penjelasan:

```go
var serverCmd = &cobra.Command{
    Use: "server",
}
```

Ini membuat subcommand bernama `server`.

Bagian ini:

```go
func init() {
    rootCmd.AddCommand(serverCmd)
}
```

mendaftarkan `serverCmd` sebagai subcommand dari `rootCmd`.

Jadi nanti command yang bisa dijalankan adalah:

```bash
go run main.go server
```

Atau kalau sudah di-build:

```bash
./app server
```

---

## cmd/worker.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
    Use:   "worker",
    Short: "Running a worker",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Worker is running!")
    },
}

func init() {
    rootCmd.AddCommand(workerCmd)
}
```

Penjelasan:

```go
var workerCmd = &cobra.Command{
    Use: "worker",
}
```

Ini membuat subcommand bernama `worker`.

Bagian ini:

```go
func init() {
    rootCmd.AddCommand(workerCmd)
}
```

mendaftarkan `workerCmd` sebagai subcommand dari `rootCmd`.

Jadi nanti command yang bisa dijalankan adalah:

```bash
go run main.go worker
```

Atau kalau sudah di-build:

```bash
./app worker
```

---

## Alur Jalannya Program

Misalnya kita menjalankan:

```bash
go run main.go server
```

Alurnya:

1. Program mulai dari `main.go`.
2. Package `cmd` di-import.
3. Go menjalankan semua `init()` di package `cmd`.
4. `serverCmd` dan `workerCmd` didaftarkan ke `rootCmd`.
5. `main()` memanggil `cmd.Execute()`.
6. `rootCmd.Execute()` membaca command dari terminal.
7. Karena command-nya `server`, maka `serverCmd.Run` dijalankan.
8. Output yang muncul:

```text
Server is running!
```

---

## Kenapa Subcommand Dipisah ke File Sendiri?

Kalau command masih sedikit, semua bisa ditaruh di satu file.

Tapi kalau command mulai banyak, lebih rapi kalau dipisah:

```text
cmd/
├── root.go
├── server.go
├── worker.go
├── migrate.go
└── user.go
```

Keuntungannya:

- Code lebih rapi.
- Setiap command punya file sendiri.
- Lebih gampang dicari.
- Lebih gampang dikembangkan.
- `main.go` tetap bersih.
- `root.go` fokus ke root command.
- File subcommand fokus ke logic command masing-masing.

---

## Contoh Lengkap

### main.go

```go
package main

import "app/cmd"

func main() {
    cmd.Execute()
}
```

### cmd/root.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

### cmd/server.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running!")
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
}
```

### cmd/worker.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
    Use:   "worker",
    Short: "Running a worker",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Worker is running!")
    },
}

func init() {
    rootCmd.AddCommand(workerCmd)
}
```

---

## Command yang Bisa Dijalankan

```bash
go run main.go
```

Output:

```text
Hello world!
```

```bash
go run main.go server
```

Output:

```text
Server is running!
```

```bash
go run main.go worker
```

Output:

```text
Worker is running!
```

---

## Intinya

Untuk project Cobra yang rapi:

- `main.go` hanya memanggil `cmd.Execute()`.
- `cmd/root.go` berisi `rootCmd` dan function `Execute()`.
- Setiap subcommand dibuat di file sendiri.
- Subcommand didaftarkan ke `rootCmd` lewat `func init()`.
- `func init()` otomatis dijalankan oleh Go sebelum `main()`.

Pattern umumnya:

```go
func init() {
    rootCmd.AddCommand(namaSubcommand)
}
```

Jadi kalau nanti mau tambah command baru, tinggal buat file baru di folder `cmd`, lalu daftarkan command-nya di `init()`.

## Cobra Flags

Di Cobra, kita bisa menambahkan **flags** ke command.

Flag adalah parameter tambahan yang dikirim lewat terminal.

Contoh:

```bash
go run main.go server --address 127.0.0.1:8080
```

Atau versi pendeknya:

```bash
go run main.go server -a 127.0.0.1:8080
```

Flag berguna supaya value tertentu bisa diatur dari command line, tanpa hardcode di dalam program.

---

## StringVarP

Untuk membuat flag bertipe string, bisa pakai:

```go
cmd.Flags().StringVarP()
```

Formatnya:

```go
cmd.Flags().StringVarP(
    &variable,
    "flag-name",
    "short-name",
    "default-value",
    "description",
)
```

Contoh:

```go
var address string

serverCmd.Flags().StringVarP(
    &address,
    "address",
    "a",
    "127.0.0.1:8080",
    "Address for server to listen to",
)
```

Penjelasan:

```go
&address
```

Variable yang akan diisi value dari flag.

```go
"address"
```

Nama flag panjang.

Dipakai seperti ini:

```bash
--address 127.0.0.1:8080
```

```go
"a"
```

Nama flag pendek.

Dipakai seperti ini:

```bash
-a 127.0.0.1:8080
```

```go
"127.0.0.1:8080"
```

Default value kalau user tidak mengisi flag.

```go
"Address for server to listen to"
```

Deskripsi flag yang akan muncul saat menjalankan help.

---

## Contoh Pemakaian Flag

Misalnya kita punya command `server`.

```go
var address string

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running on", address)
    },
}

func init() {
    serverCmd.Flags().StringVarP(
        &address,
        "address",
        "a",
        "127.0.0.1:8080",
        "Address for server to listen to",
    )

    rootCmd.AddCommand(serverCmd)
}
```

Kalau dijalankan tanpa flag:

```bash
go run main.go server
```

Output:

```text
Server is running on 127.0.0.1:8080
```

Kalau dijalankan dengan flag:

```bash
go run main.go server -a 0.0.0.0:8080
```

Output:

```text
Server is running on 0.0.0.0:8080
```

Atau pakai flag panjang:

```bash
go run main.go server --address 0.0.0.0:8080
```

Output:

```text
Server is running on 0.0.0.0:8080
```

---

## Catatan Address

Kalau menggunakan:

```text
127.0.0.1:8080
```

Server hanya bisa diakses dari local machine.

Kalau menggunakan:

```text
0.0.0.0:8080
```

Server listen ke semua network interface, jadi bisa diakses dari luar machine kalau firewall dan port-nya terbuka.

Bisa juga ditulis:

```text
:8080
```

Ini biasanya berarti listen di semua interface pada port `8080`.

Yang salah:

```text
:0.0.0.0:8080
```

Karena format address-nya tidak valid.

---

## Help Command

Setelah flag ditambahkan, Cobra otomatis menampilkan flag tersebut di help command.

Contoh:

```bash
go run main.go server -h
```

Atau:

```bash
go run main.go server --help
```

Nanti akan muncul keterangan seperti:

```text
Running a server

Usage:
  app server [flags]

Flags:
  -a, --address string   Address for server to listen to (default "127.0.0.1:8080")
  -h, --help             help for server
```

Jadi dengan menambahkan:

```go
serverCmd.Flags().StringVarP(
    &address,
    "address",
    "a",
    "127.0.0.1:8080",
    "Address for server to listen to",
)
```

Cobra otomatis menambahkan flag ke command `server` dan otomatis menampilkannya di `-h`.

---

# Contoh Struktur Cobra dengan Flags

Best practice untuk project Cobra adalah memisahkan command ke beberapa file.

Struktur project:

```text
project/
├── main.go
└── cmd/
    ├── root.go
    ├── server.go
    └── worker.go
```

---

## main.go

```go
package main

import "app/cmd"

func main() {
    cmd.Execute()
}
```

`main.go` cukup memanggil:

```go
cmd.Execute()
```

Jadi logic command tidak ditaruh semua di `main.go`.

---

## cmd/root.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Penjelasan:

```go
var rootCmd = &cobra.Command{
    Use: "app",
}
```

Ini adalah command utama.

```go
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

Function `Execute()` dibuat supaya dari `main.go`, kita cukup memanggil:

```go
cmd.Execute()
```

---

## cmd/server.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var address string

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running on", address)
    },
}

func init() {
    serverCmd.Flags().StringVarP(
        &address,
        "address",
        "a",
        "127.0.0.1:8080",
        "Address for server to listen to",
    )

    rootCmd.AddCommand(serverCmd)
}
```

Penjelasan:

```go
var address string
```

Variable ini akan menyimpan value dari flag `--address` atau `-a`.

```go
serverCmd.Flags().StringVarP(
    &address,
    "address",
    "a",
    "127.0.0.1:8080",
    "Address for server to listen to",
)
```

Bagian ini mendaftarkan flag ke `serverCmd`.

```go
rootCmd.AddCommand(serverCmd)
```

Bagian ini mendaftarkan `serverCmd` sebagai subcommand dari `rootCmd`.

Jadi command yang bisa dijalankan:

```bash
go run main.go server
```

```bash
go run main.go server -a 0.0.0.0:8080
```

```bash
go run main.go server --address 0.0.0.0:8080
```

---

## cmd/worker.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
    Use:   "worker",
    Short: "Running a worker",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Worker is running!")
    },
}

func init() {
    rootCmd.AddCommand(workerCmd)
}
```

Command ini bisa dijalankan dengan:

```bash
go run main.go worker
```

Output:

```text
Worker is running!
```

---

# func init()

Di Go ada function khusus bernama `init`.

```go
func init() {
    // initialization logic
}
```

`init()` akan otomatis dijalankan oleh Go ketika package di-load.

Catatan:

- `init()` tidak perlu dipanggil manual.
- `init()` otomatis dijalankan sebelum `main()`.
- `init()` biasanya dipakai untuk setup awal.
- Dalam Cobra, `init()` sering dipakai untuk mendaftarkan subcommand dan flags.
- `init()` hanya dijalankan sekali dalam satu aplikasi berjalan.

Contoh:

```go
func init() {
    rootCmd.AddCommand(serverCmd)
}
```

Atau kalau command punya flag:

```go
func init() {
    serverCmd.Flags().StringVarP(
        &address,
        "address",
        "a",
        "127.0.0.1:8080",
        "Address for server to listen to",
    )

    rootCmd.AddCommand(serverCmd)
}
```

---

# Alur Jalannya Program

Misalnya menjalankan:

```bash
go run main.go server -a 0.0.0.0:8080
```

Alurnya:

1. Program mulai dari `main.go`.
2. `main.go` meng-import package `cmd`.
3. Go menjalankan semua `init()` di package `cmd`.
4. Di `server.go`, flag `address` didaftarkan ke `serverCmd`.
5. `serverCmd` didaftarkan ke `rootCmd`.
6. `main()` memanggil `cmd.Execute()`.
7. `rootCmd.Execute()` membaca command dari terminal.
8. Karena command-nya `server`, maka `serverCmd.Run` dijalankan.
9. Karena ada flag `-a 0.0.0.0:8080`, variable `address` berisi `"0.0.0.0:8080"`.
10. Output yang muncul:

```text
Server is running on 0.0.0.0:8080
```

---

# Local Flags vs Persistent Flags

Kalau pakai:

```go
serverCmd.Flags()
```

Maka flag hanya berlaku untuk command tersebut.

Contoh:

```go
serverCmd.Flags().StringVarP(
    &address,
    "address",
    "a",
    "127.0.0.1:8080",
    "Address for server to listen to",
)
```

Flag `address` hanya berlaku untuk:

```bash
go run main.go server -a 127.0.0.1:8080
```

Kalau pakai:

```go
rootCmd.PersistentFlags()
```

Maka flag berlaku untuk root command dan subcommand di bawahnya.

Contoh:

```go
rootCmd.PersistentFlags().StringVarP(
    &configPath,
    "config",
    "c",
    "config.yaml",
    "Path to config file",
)
```

Flag ini bisa dipakai di banyak command:

```bash
go run main.go --config config.yaml
go run main.go server --config config.yaml
go run main.go worker --config config.yaml
```

---

# Intinya

Untuk menambahkan flag di Cobra, gunakan:

```go
cmd.Flags().StringVarP(
    &variable,
    "flag-name",
    "short-name",
    "default-value",
    "description",
)
```

Contoh:

```go
serverCmd.Flags().StringVarP(
    &address,
    "address",
    "a",
    "127.0.0.1:8080",
    "Address for server to listen to",
)
```

Dengan ini, command bisa dijalankan seperti:

```bash
go run main.go server -a 0.0.0.0:8080
```

Atau:

```bash
go run main.go server --address 0.0.0.0:8080
```

Cobra juga otomatis membuat help command:

```bash
go run main.go server -h
```

Yang akan menampilkan usage dan daftar flags.

## Persistent Flags di Cobra

Selain `Flags()`, Cobra juga punya `PersistentFlags()`.

Kalau memakai:

```go
rootCmd.PersistentFlags().StringVarP(
    &config,
    "config",
    "c",
    "config.yaml",
    "Config file to use",
)
```

artinya kita membuat flag `config` yang berlaku untuk `rootCmd` dan semua subcommand di bawahnya.

---

## Contoh

```go
var config string

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Using config:", config)
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(
        &config,
        "config",
        "c",
        "config.yaml",
        "Config file to use",
    )
}
```

Penjelasan:

- `&config` adalah variable yang akan menyimpan value dari flag.
- `"config"` adalah nama flag panjang, dipakai dengan `--config`.
- `"c"` adalah nama flag pendek, dipakai dengan `-c`.
- `"config.yaml"` adalah default value.
- `"Config file to use"` adalah deskripsi yang muncul di help.

---

## Cara Menjalankan

```bash
go run main.go --config dev.yaml
```

Atau:

```bash
go run main.go -c dev.yaml
```

Karena ini `PersistentFlags`, flag ini juga bisa dipakai di subcommand:

```bash
go run main.go server --config dev.yaml
```

```bash
go run main.go worker -c dev.yaml
```

---

## Flags vs PersistentFlags

Kalau pakai:

```go
serverCmd.Flags()
```

flag hanya berlaku untuk command `server`.

Kalau pakai:

```go
rootCmd.PersistentFlags()
```

flag berlaku untuk `rootCmd` dan semua subcommand.

Jadi biasanya:

- `Flags()` dipakai untuk flag khusus satu command.
- `PersistentFlags()` dipakai untuk flag global, seperti config file, environment, atau debug mode.


## PersistentPreRun di Cobra

`PersistentPreRun` adalah function yang dijalankan **sebelum `Run` command dijalankan**.

Kalau `PersistentPreRun` dipasang di `rootCmd`, maka dia bisa jalan sebelum command utama maupun subcommand.

Biasanya ini dipakai untuk setup awal, misalnya:

- load config file
- connect database
- setup logger
- validasi environment
- membaca global flag seperti `--config`

---

## Contoh PersistentPreRun

```go
var config string

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",

    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("Loading config from:", config)

        // contoh:
        // loadConfig(config)
        // setupLogger()
        // connectDatabase()
    },

    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(
        &config,
        "config",
        "c",
        "config.yaml",
        "Config file to use",
    )
}
```

---

## Contoh Subcommand

```go
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running!")
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
}
```

Kalau dijalankan:

```bash
go run main.go server --config dev.yaml
```

Maka alurnya:

```text
PersistentPreRun rootCmd jalan dulu
Run serverCmd jalan
```

Output kira-kira:

```text
Loading config from: dev.yaml
Server is running!
```

---

## PersistentPreRun vs Run

```go
PersistentPreRun: func(cmd *cobra.Command, args []string) {
    // dijalankan sebelum Run
},
Run: func(cmd *cobra.Command, args []string) {
    // logic utama command
},
```

`PersistentPreRun` cocok untuk logic yang harus jalan sebelum command utama.

`Run` cocok untuk logic utama dari command tersebut.

---

## Kenapa Berguna?

Misalnya semua command butuh config.

```bash
go run main.go server --config dev.yaml
go run main.go worker --config dev.yaml
go run main.go migrate --config dev.yaml
```

Daripada load config di setiap command satu-satu, kita bisa taruh di `PersistentPreRun` root.

Jadi semua subcommand bisa otomatis menjalankan setup yang sama sebelum logic command-nya berjalan.

---

## Intinya

- `PersistentPreRun` jalan sebelum `Run`.
- Kalau dipasang di `rootCmd`, bisa dipakai oleh subcommand.
- Cocok untuk setup global seperti load config.
- Biasanya dipakai bareng `PersistentFlags`.
- Contoh umum: `--config` dibaca dari flag, lalu config di-load di `PersistentPreRun`.

## Load Config dengan Viper di PersistentPreRun

Misalnya kita punya file config:

```yaml
# config.yaml
server:
  address: "127.0.0.1:8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "app_db"
```

Lalu di Cobra, kita bisa pakai `PersistentPreRun` atau `PersistentPreRunE` untuk load config sebelum command dijalankan.

Biasanya lebih enak pakai `PersistentPreRunE`, karena bisa return error kalau config gagal dibaca.

---

## Contoh root.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var config string

type AppConfig struct {
    Server struct {
        Address string `mapstructure:"address"`
    } `mapstructure:"server"`

    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        User     string `mapstructure:"user"`
        Password string `mapstructure:"password"`
        Name     string `mapstructure:"name"`
    } `mapstructure:"database"`
}

var appConfig AppConfig

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Main application command",

    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        viper.SetConfigFile(config)

        err := viper.ReadInConfig()
        if err != nil {
            return fmt.Errorf("failed to read config: %w", err)
        }

        err = viper.Unmarshal(&appConfig)
        if err != nil {
            return fmt.Errorf("failed to parse config: %w", err)
        }

        fmt.Println("Using config file:", viper.ConfigFileUsed())

        return nil
    },

    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello world!")
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(
        &config,
        "config",
        "c",
        "config.yaml",
        "Config file to use",
    )
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }
}
```

---

## Contoh server.go

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Running a server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Server is running on", appConfig.Server.Address)
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
}
```

---

## Cara Menjalankan

```bash
go run main.go server
```

Secara default akan membaca:

```bash
config.yaml
```

Kalau mau pakai config lain:

```bash
go run main.go server --config dev.yaml
```

Atau versi pendek:

```bash
go run main.go server -c dev.yaml
```

---

## Penjelasan

Bagian ini:

```go
viper.SetConfigFile(config)
```

digunakan untuk menentukan file config yang mau dibaca.

Bagian ini:

```go
err := viper.ReadInConfig()
```

digunakan untuk membaca isi file config.

Bagian ini:

```go
err = viper.Unmarshal(&appConfig)
```

digunakan untuk mengubah isi config menjadi struct Go.

Jadi kalau di `config.yaml` ada:

```yaml
server:
  address: "127.0.0.1:8080"
```

maka di Go bisa diakses lewat:

```go
appConfig.Server.Address
```

---

## Kenapa Ditaruh di PersistentPreRun?

Karena config biasanya dibutuhkan oleh banyak command.

Misalnya:

```bash
go run main.go server
go run main.go worker
go run main.go migrate
```

Semua command itu mungkin butuh config.

Daripada load config di setiap command satu-satu, lebih rapi kalau load config ditaruh di `PersistentPreRunE` milik `rootCmd`.

Jadi sebelum subcommand seperti `server`, `worker`, atau `migrate` jalan, config sudah dibaca dulu.

---

## Intinya

Pattern-nya:

```go
rootCmd.PersistentFlags().StringVarP(
    &config,
    "config",
    "c",
    "config.yaml",
    "Config file to use",
)
```

Lalu di `PersistentPreRunE`:

```go
viper.SetConfigFile(config)

err := viper.ReadInConfig()
if err != nil {
    return err
}

err = viper.Unmarshal(&appConfig)
if err != nil {
    return err
}
```

Dengan ini, config dari file YAML bisa dibaca sekali di awal, lalu dipakai oleh semua subcommand.

# Database di Go

Di Go, database biasanya diakses lewat package:

```go
database/sql
```

Package `database/sql` adalah interface umum dari Go untuk berkomunikasi dengan database.

Go tidak menyediakan implementasi driver database secara langsung. Jadi untuk database tertentu, kita tetap perlu install driver external.

Contoh driver:

```text
SQLite     -> github.com/mattn/go-sqlite3
PostgreSQL -> github.com/lib/pq
SQL Server -> github.com/microsoft/go-mssqldb
```

Konsepnya:

```text
database/sql = interface umum dari Go
driver       = implementasi spesifik untuk database tertentu
```

Jadi walaupun databasenya beda-beda, cara pakainya di Go kurang lebih mirip karena sama-sama lewat `database/sql`.

---

## SQL Interface vs ORM

Kalau pakai `database/sql`, kita menulis query SQL secara manual.

Kelebihan:

- Lebih fleksibel untuk query complex.
- Lebih dekat dengan SQL asli.
- Tidak terlalu banyak abstraction.
- Cocok kalau ingin kontrol penuh terhadap query.

Kekurangan:

- Kode lebih panjang.
- Harus handle `Scan`, `rows.Close`, error, dan query manual.
- Ada sedikit overhead karena memakai interface Go.

Kalau ingin ORM, bisa pakai library seperti:

```text
gorm
```

ORM cocok untuk query simple atau CRUD biasa. Tapi untuk query complex, ORM kadang kurang fleksibel dan bisa punya bottleneck sendiri.

---

# SQLite di Go

Untuk SQLite, install driver:

```bash
go get github.com/mattn/go-sqlite3
```

Import driver-nya:

```go
import _ "github.com/mattn/go-sqlite3"
```

Pakai `_` karena driver ini tidak dipakai langsung di kode, tapi hanya didaftarkan ke `database/sql`.

---

## Membuka Database

```go
db, err := sql.Open("sqlite3", "file:db.sqlite3")
if err != nil {
    log.Println(err)
    return
}
defer db.Close()
```

Penjelasan:

```go
sql.Open("sqlite3", "file:db.sqlite3")
```

- `"sqlite3"` adalah nama driver.
- `"file:db.sqlite3"` adalah lokasi file database SQLite.

Catatan:

`sql.Open()` belum benar-benar mengecek koneksi database. Untuk memastikan koneksi bisa dipakai, gunakan:

```go
err = db.PingContext(ctx)
if err != nil {
    log.Println(err)
    return
}
```

---

# Exec, Query, dan QueryRow

Di `database/sql`, ada 3 method utama yang sering dipakai.

## ExecContext

Dipakai untuk query yang tidak mengembalikan row.

Biasanya untuk:

```text
CREATE
INSERT
UPDATE
DELETE
```

Contoh:

```go
result, err := db.ExecContext(ctx, query)
```

Dari `result`, kita bisa mengambil:

```go
lastID, err := result.LastInsertId()
```

atau:

```go
affected, err := result.RowsAffected()
```

Catatan:

`LastInsertId()` dan `RowsAffected()` tergantung implementasi driver. Tidak semua driver mendukung dua method ini dengan cara yang sama.

---

## QueryContext

Dipakai untuk query yang mengembalikan banyak row.

Contoh:

```go
rows, err := db.QueryContext(ctx, "SELECT id, name, age FROM users")
```

Karena hasilnya banyak row, kita perlu loop:

```go
for rows.Next() {
    // scan data
}
```

---

## QueryRowContext

Dipakai untuk query yang hanya mengembalikan satu row.

Contoh:

```go
row := db.QueryRowContext(ctx, "SELECT id, name, age FROM users WHERE id = ?", id)
```

Lalu datanya diambil dengan:

```go
err := row.Scan(&u.ID, &u.Name, &u.Age)
```

---

# Context

`context.Context` digunakan untuk membawa informasi seperti:

- timeout
- cancel signal
- deadline
- request scope

Dalam database operation, context berguna supaya query bisa dibatalkan kalau terlalu lama.

Contoh context biasa:

```go
ctx := context.Background()
```

Contoh context dengan timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

Lalu context dipakai di query:

```go
db.ExecContext(ctx, query)
db.QueryContext(ctx, query)
db.QueryRowContext(ctx, query)
```

Jadi kalau query terlalu lama, context bisa menghentikan query tersebut.

---

# Repo Pattern

Biasanya logic database tidak langsung ditaruh semua di `main.go`.

Kita bisa membuat repository.

```go
type Repo struct {
    db *sql.DB
}
```

Gunanya:

- Menyimpan koneksi database.
- Mengelompokkan query yang berhubungan.
- Membuat kode lebih rapi.
- Method seperti `Insert`, `List`, dan `Get` tidak perlu menerima `db` terus-menerus.

Constructor-nya:

```go
func NewRepo(db *sql.DB) *Repo {
    return &Repo{
        db: db,
    }
}
```

Kenapa return `*Repo`?

Karena repository biasanya tidak perlu di-copy. Kita cukup pakai pointer ke struct yang menyimpan koneksi database.

---

# Contoh File users/db.go

```go
package users

import (
    "context"
    "database/sql"
    "errors"
)

type User struct {
    ID   int
    Name string
    Age  int
}

type Repo struct {
    db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
    return &Repo{
        db: db,
    }
}

const createTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL
);
`

func (r *Repo) CreateTable(ctx context.Context) error {
    _, err := r.db.ExecContext(ctx, createTableQuery)
    return err
}

func (r *Repo) Insert(ctx context.Context, name string, age int) error {
    query := `
    INSERT INTO users (name, age)
    VALUES (?, ?);
    `

    _, err := r.db.ExecContext(ctx, query, name, age)
    return err
}

func (r *Repo) List(ctx context.Context) ([]User, error) {
    query := `
    SELECT id, name, age
    FROM users;
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []User{}

    for rows.Next() {
        var u User

        err := rows.Scan(&u.ID, &u.Name, &u.Age)
        if err != nil {
            return nil, err
        }

        users = append(users, u)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

func (r *Repo) Get(ctx context.Context, id int) (User, error) {
    query := `
    SELECT id, name, age
    FROM users
    WHERE id = ?;
    `

    var u User

    row := r.db.QueryRowContext(ctx, query, id)

    err := row.Scan(&u.ID, &u.Name, &u.Age)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return User{}, err
        }

        return User{}, err
    }

    return u, nil
}
```

---

# Kenapa Jangan Pakai fmt.Sprintf untuk Query?

Jangan membuat query seperti ini:

```go
query := fmt.Sprintf(
    "INSERT INTO users (name, age) VALUES ('%s', %d)",
    name,
    age,
)
```

Karena bisa terkena SQL injection.

Yang benar adalah pakai parameter:

```go
_, err := r.db.ExecContext(
    ctx,
    "INSERT INTO users (name, age) VALUES (?, ?)",
    name,
    age,
)
```

Dengan cara ini, value `name` dan `age` dikirim sebagai parameter, bukan digabung langsung ke string SQL.

Catatan:

Untuk SQLite, placeholder yang digunakan adalah:

```sql
?
```

Untuk PostgreSQL, placeholder biasanya:

```sql
$1, $2, $3
```

---

# Transaction

Transaction digunakan kalau kita ingin beberapa perubahan database dianggap sebagai satu kesatuan.

Misalnya ada 3 query:

```text
INSERT user
INSERT profile
INSERT log
```

Kalau salah satu gagal, semua perubahan sebelumnya harus dibatalkan.

Di transaction ada dua operasi penting:

```go
tx.Commit()
```

Digunakan untuk menyimpan semua perubahan.

```go
tx.Rollback()
```

Digunakan untuk membatalkan semua perubahan.

Setelah `BeginTx`, transaction harus selalu diakhiri dengan salah satu:

```go
Commit()
```

atau:

```go
Rollback()
```

Kalau tidak, koneksi ke database bisa menggantung dan bikin masalah.

---

## Contoh Transaction

```go
func (r *Repo) InsertManyTx(ctx context.Context, users []User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer tx.Rollback()

    query := `
    INSERT INTO users (name, age)
    VALUES (?, ?);
    `

    for _, user := range users {
        _, err := tx.ExecContext(ctx, query, user.Name, user.Age)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

Penjelasan:

```go
tx, err := r.db.BeginTx(ctx, nil)
```

Digunakan untuk memulai transaction.

```go
defer tx.Rollback()
```

Digunakan sebagai safety net.

Kalau function return lebih awal karena error, transaction akan otomatis di-rollback.

```go
return tx.Commit()
```

Kalau semua query berhasil, transaction di-commit.

Setelah `Commit()` berhasil, `defer tx.Rollback()` tetap akan terpanggil, tapi rollback akan error karena transaction sudah selesai. Biasanya error ini diabaikan.

---

## Transaction dengan Rollback Manual

Kalau ingin lebih eksplisit, bisa juga seperti ini:

```go
func (r *Repo) InsertManyTx(ctx context.Context, users []User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    query := `
    INSERT INTO users (name, age)
    VALUES (?, ?);
    `

    for _, user := range users {
        _, err := tx.ExecContext(ctx, query, user.Name, user.Age)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    err = tx.Commit()
    if err != nil {
        tx.Rollback()
        return err
    }

    return nil
}
```

Versi ini lebih panjang, tapi lebih eksplisit.

Untuk kode yang sederhana, pattern ini sudah cukup aman:

```go
tx, err := r.db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// query-query

return tx.Commit()
```

---

# Contoh main.go

```go
package main

import (
    "context"
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"

    "your-module/users"
)

func main() {
    ctx := context.Background()

    db, err := sql.Open("sqlite3", "file:db.sqlite3")
    if err != nil {
        log.Println(err)
        return
    }
    defer db.Close()

    err = db.PingContext(ctx)
    if err != nil {
        log.Println(err)
        return
    }

    userRepo := users.NewRepo(db)

    err = userRepo.CreateTable(ctx)
    if err != nil {
        log.Println(err)
        return
    }

    err = userRepo.Insert(ctx, "John", 17)
    if err != nil {
        log.Println(err)
        return
    }

    err = userRepo.Insert(ctx, "Jane", 18)
    if err != nil {
        log.Println(err)
        return
    }

    userList, err := userRepo.List(ctx)
    if err != nil {
        log.Println(err)
        return
    }

    log.Println(userList)

    user, err := userRepo.Get(ctx, 1)
    if err != nil {
        log.Println(err)
        return
    }

    log.Println(user)

    err = userRepo.InsertManyTx(ctx, []users.User{
        {
            Name: "Alice",
            Age:  20,
        },
        {
            Name: "Bob",
            Age:  21,
        },
    })
    if err != nil {
        log.Println(err)
        return
    }
}
```

---

# InitTable Function

Kalau mau logic init database dipisah ke function sendiri, bisa seperti ini:

```go
func InitTable(ctx context.Context, db *sql.DB) error {
    userRepo := users.NewRepo(db)

    err := userRepo.CreateTable(ctx)
    if err != nil {
        return err
    }

    err = userRepo.Insert(ctx, "John", 17)
    if err != nil {
        return err
    }

    err = userRepo.Insert(ctx, "Jane", 18)
    if err != nil {
        return err
    }

    return nil
}
```

Pemakaian:

```go
ctx := context.Background()

err := InitTable(ctx, db)
if err != nil {
    log.Println(err)
    return
}
```

---

# Intinya

Go memakai package `database/sql` sebagai interface umum untuk database.

Untuk SQLite, kita perlu driver:

```go
_ "github.com/mattn/go-sqlite3"
```

Untuk membuka database:

```go
db, err := sql.Open("sqlite3", "file:db.sqlite3")
```

Untuk query:

```go
db.ExecContext(ctx, query)
db.QueryContext(ctx, query)
db.QueryRowContext(ctx, query)
```

Gunakan:

- `ExecContext` untuk query yang tidak mengembalikan row.
- `QueryContext` untuk query yang mengembalikan banyak row.
- `QueryRowContext` untuk query yang mengembalikan satu row.

Gunakan repository pattern supaya kode database lebih rapi:

```go
type Repo struct {
    db *sql.DB
}
```

Gunakan parameter query untuk mencegah SQL injection:

```go
r.db.ExecContext(ctx, "INSERT INTO users (name, age) VALUES (?, ?)", name, age)
```

Gunakan transaction kalau beberapa query harus dianggap sebagai satu kesatuan:

```go
tx, err := r.db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// query-query

return tx.Commit()
```

# JSON

Di Go, JSON biasanya dipakai lewat package:

```go
encoding/json
```

JSON sering digunakan untuk mengubah data Go menjadi format JSON, atau sebaliknya.

---

## Struct Tag

Kita bisa memberi instruksi khusus ke field struct menggunakan tag.

Contoh:

```go
type User struct {
    Name     string `json:"name" db:"name" yaml:"name"`
    Age      int    `json:"age"`
    Address  string `json:"address,omitempty"`
    Password string `json:"-"`
}
```

Penjelasan:

```go
Name string `json:"name"`
```

Artinya saat diubah menjadi JSON, field `Name` akan muncul sebagai `name`.

Contoh output:

```json
{
  "name": "John"
}
```

Field struct harus diawali huruf besar supaya bisa dibaca oleh package `encoding/json`.

Contoh yang bisa dibaca:

```go
Name string
Age  int
```

Contoh yang tidak bisa dibaca:

```go
name string
age  int
```

Karena field dengan huruf kecil bersifat private di package Go.

---

## json:"-"

Kalau field tidak mau ditampilkan di JSON, gunakan:

```go
Password string `json:"-"`
```

Contoh:

```go
type User struct {
    Name     string `json:"name"`
    Age      int    `json:"age"`
    Password string `json:"-"`
}
```

Kalau di-marshal, `Password` tidak akan muncul di JSON.

---

## omitempty

`omitempty` digunakan supaya field tidak muncul kalau nilainya kosong atau zero value.

Contoh:

```go
type User struct {
    Name    string `json:"name"`
    Address string `json:"address,omitempty"`
}
```

Kalau `Address` kosong, output JSON tidak akan menampilkan field `address`.

Zero value yang dianggap kosong:

```text
string  -> ""
int     -> 0
bool    -> false
slice   -> nil
map     -> nil
pointer -> nil
```

---

## Encoding: Struct Go ke JSON

Encoding berarti mengubah data Go menjadi JSON.

Contoh dengan `json.Marshal`:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type User struct {
    Name     string `json:"name"`
    Age      int    `json:"age"`
    Address  string `json:"address,omitempty"`
    Password string `json:"-"`
}

func main() {
    u := User{
        Name:     "John",
        Age:      17,
        Address:  "BSD",
        Password: "secret",
    }

    b, err := json.Marshal(u)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(string(b))
}
```

Output:

```json
{"name":"John","age":17,"address":"BSD"}
```

Field `Password` tidak muncul karena memakai:

```go
json:"-"
```

---

## Decoding: JSON ke Struct Go

Decoding berarti mengubah JSON menjadi data Go.

Contoh dengan `json.Unmarshal`:

```go
func DecodeExample() {
    str := `{"name":"John","age":18}`

    var u User

    err := json.Unmarshal([]byte(str), &u)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(u.Name, u.Age)
}
```

Penjelasan:

```go
json.Unmarshal([]byte(str), &u)
```

- `[]byte(str)` mengubah string menjadi byte.
- `&u` berarti hasil parsing JSON akan dimasukkan ke variable `u`.

Kalau JSON punya field yang tidak ada di struct, field itu akan diabaikan.

Contoh:

```json
{"first_name":"John","age":18}
```

Kalau struct-nya hanya punya:

```go
Name string `json:"name"`
Age  int    `json:"age"`
```

maka `first_name` tidak masuk ke `Name`, karena tag-nya beda.

Akibatnya `Name` tetap berisi zero value, yaitu string kosong.

---

## JSON Command dengan Cobra

Contoh command Cobra untuk output JSON:

```go
func init() {
    rootCmd.AddCommand(&cobra.Command{
        Use:   "json",
        Short: "Print user as JSON",
        Run: func(cmd *cobra.Command, args []string) {
            u := User{
                Name: "John",
                Age:  17,
            }

            b, err := json.Marshal(u)
            if err != nil {
                log.Println(err)
                return
            }

            fmt.Println(string(b))
        },
    })
}
```

Kalau dijalankan:

```bash
go run main.go json
```

Output:

```json
{"name":"John","age":17}
```

---

# HTTP Server

Di Go, HTTP server bawaan ada di package:

```go
net/http
```

Kita bisa membuat server tanpa framework tambahan.

---

## Contoh HTTP Server Bawaan

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func StartServer(address string) {
    handler := http.NewServeMux()

    handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello world"))
    })

    s := http.Server{
        Addr:    address,
        Handler: handler,
    }

    err := s.ListenAndServe()
    if err != nil {
        log.Println(err)
    }
}
```

Penjelasan:

```go
handler := http.NewServeMux()
```

`ServeMux` adalah router bawaan Go. Dia menentukan path mana yang akan memanggil handler tertentu.

```go
handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    // logic
})
```

Artinya kalau user membuka path `/hello`, function ini akan dijalankan.

---

## ResponseWriter dan Request

```go
func(w http.ResponseWriter, r *http.Request)
```

Penjelasan:

```go
w http.ResponseWriter
```

Digunakan untuk membuat response ke client.

Contoh yang diatur lewat `ResponseWriter`:

- status code
- response header
- response body

Contoh:

```go
w.WriteHeader(http.StatusOK)
w.Write([]byte("Hello world"))
```

```go
r *http.Request
```

Berisi informasi request dari client.

Contoh isi request:

- method
- path
- header
- body
- query parameter

Contoh cek method:

```go
if r.Method == http.MethodPost {
    // handle POST
}
```

---

## http.Server

```go
s := http.Server{
    Addr:    address,
    Handler: handler,
}
```

Penjelasan:

```go
Addr: address
```

Alamat server, misalnya:

```text
127.0.0.1:8080
```

atau:

```text
0.0.0.0:8080
```

```go
Handler: handler
```

Menentukan router atau handler utama yang akan dipakai oleh server.

```go
s.ListenAndServe()
```

Menjalankan server.

---

# Gin Framework

Gin adalah framework Go untuk membuat HTTP server.

Install:

```bash
go get github.com/gin-gonic/gin
```

---

## Contoh Gin Server

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func StartGinServer(address string) {
    handler := gin.New()

    handler.Use(func(ctx *gin.Context) {
        log.Println(ctx.Request.Method, ctx.Request.URL.Path)

        ctx.Next()
    })

    handler.GET("/hello", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "Hello world!")
    })

    s := http.Server{
        Addr:    address,
        Handler: handler,
    }

    err := s.ListenAndServe()
    if err != nil {
        log.Println(err)
    }
}
```

Penjelasan:

```go
handler := gin.New()
```

Membuat Gin engine/router.

```go
handler.Use(...)
```

Menambahkan middleware.

```go
handler.GET("/hello", ...)
```

Mendaftarkan route GET `/hello`.

```go
ctx.String(http.StatusOK, "Hello world!")
```

Mengirim response text dengan status `200 OK`.

Kalau request ke path yang tidak ada, misalnya `/`, hasilnya akan:

```text
404 page not found
```

Kalau request ke `/hello`, hasilnya:

```text
Hello world!
```

---

## gin.Context vs context.Context

`gin.Context` dan `context.Context` itu beda.

```go
ctx *gin.Context
```

Dipakai oleh Gin untuk handle request dan response.

Di dalamnya ada helper seperti:

```go
ctx.String()
ctx.JSON()
ctx.Request
ctx.Writer
```

Sedangkan:

```go
context.Context
```

adalah context standar Go, biasanya dipakai untuk timeout, cancel, deadline, dan request scope.

Kalau butuh `context.Context` dari Gin, ambil dari request:

```go
goCtx := ctx.Request.Context()
```

Biasanya ini dipakai kalau mau query database:

```go
db.QueryContext(goCtx, query)
```

---

# Chi Framework

Chi adalah framework/router Go yang dekat dengan `net/http`.

Install:

```bash
go get github.com/go-chi/chi/v5
```

Chi banyak menggunakan pattern `http.Handler`, sehingga middleware-nya terlihat lebih eksplisit.

---

## Contoh Chi Server

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
)

func StartChiServer(address string) {
    handler := chi.NewRouter()

    handler.Use(LogMiddleware)

    handler.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Hello world")
    })

    s := http.Server{
        Addr:    address,
        Handler: handler,
    }

    err := s.ListenAndServe()
    if err != nil {
        log.Println(err)
    }
}
```

---

# Middleware

Middleware adalah logic yang dijalankan sebelum atau sesudah handler utama.

Contoh alurnya:

```text
middleware 1 -> middleware 2 -> middleware 3 -> controller
```

Middleware berguna untuk:

- logging
- authorization
- authentication
- recovery
- parsing request
- validasi header
- setup request context

Go tidak punya keyword khusus bernama middleware, tapi middleware bisa dibuat dengan pattern `http.Handler`.

---

## Middleware di Chi

Bentuk middleware di Chi:

```go
func MiddlewareName(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // logic sebelum handler utama

        next.ServeHTTP(w, r)

        // logic setelah handler utama
    })
}
```

Contoh logging middleware:

```go
func LogMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("From middleware:", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)
    })
}
```

Penjelasan:

```go
next http.Handler
```

adalah handler berikutnya yang akan dijalankan.

```go
next.ServeHTTP(w, r)
```

digunakan untuk melanjutkan request ke handler berikutnya.

Middleware ini seperti wrapper.

Middleware membungkus handler utama, lalu menentukan apakah request boleh lanjut atau tidak.

---

## Authorization Middleware

Contoh middleware untuk mengecek header `Authorization`:

```go
func AuthorizationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "cuma login yang bisa lanjut")
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

Penjelasan:

```go
r.Header.Get("Authorization")
```

Digunakan untuk mengambil value header `Authorization`.

Kalau kosong, response langsung diberi:

```go
http.StatusUnauthorized
```

Lalu `return`, supaya request tidak lanjut ke handler utama.

Kalau header ada, lanjut ke handler berikutnya:

```go
next.ServeHTTP(w, r)
```

---

## Memakai Middleware untuk Semua Route

```go
handler := chi.NewRouter()

handler.Use(AuthorizationMiddleware)

handler.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello world")
})
```

Karena pakai:

```go
handler.Use(AuthorizationMiddleware)
```

maka semua route di router tersebut akan kena middleware authorization.

---

## Group Route yang Protected

Kalau ada route public dan route protected, bisa pakai group.

```go
handler := chi.NewRouter()

handler.Get("/public", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "This is public")
})

handler.Group(func(r chi.Router) {
    r.Use(AuthorizationMiddleware)

    r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "This is protected")
    })
})
```

Penjelasan:

- `/public` tidak kena authorization middleware.
- `/profile` kena authorization middleware.

---

# Output JSON di HTTP Route

Untuk mengirim JSON response, set header:

```go
w.Header().Set("Content-Type", "application/json")
```

Contoh dengan `json.Marshal`:

```go
handler.Get("/user", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    u := User{
        Name:    "John",
        Age:     17,
        Address: "BSD",
    }

    b, err := json.Marshal(u)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "failed to encode response")
        return
    }

    w.Write(b)
})
```

Cara yang lebih ringkas adalah pakai `json.NewEncoder`.

```go
handler.Get("/user", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    u := User{
        Name:    "John",
        Age:     17,
        Address: "BSD",
    }

    json.NewEncoder(w).Encode(u)
})
```

Penjelasan:

```go
json.NewEncoder(w).Encode(u)
```

Artinya data `u` langsung diubah menjadi JSON dan ditulis ke `w`.

Karena `w` adalah `http.ResponseWriter`, maka `w` bisa menjadi target output untuk encoder.

---

# POST JSON Request

Kalau menerima JSON dari request body, gunakan:

```go
json.NewDecoder(r.Body).Decode(&payload)
```

Contoh:

```go
type CreateUserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

handler.Post("/users", func(w http.ResponseWriter, r *http.Request) {
    var payload CreateUserRequest

    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "unable to decode payload")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(payload)
})
```

Penjelasan:

```go
json.NewDecoder(r.Body)
```

Membuat decoder yang membaca JSON dari request body.

```go
Decode(&payload)
```

Mengubah JSON dari body menjadi struct Go.

Jadi kalau request body-nya:

```json
{
  "name": "John",
  "age": 17
}
```

maka hasilnya masuk ke:

```go
payload.Name
payload.Age
```

---

## Encoder vs Decoder

```go
json.NewEncoder(w).Encode(data)
```

Digunakan untuk menulis data Go menjadi JSON response.

Biasanya dipakai untuk output response.

```go
json.NewDecoder(r.Body).Decode(&data)
```

Digunakan untuk membaca JSON request body menjadi data Go.

Biasanya dipakai untuk input request.

Ringkasnya:

```text
Encoder -> Go data ke JSON
Decoder -> JSON ke Go data
```

---

# Contoh Chi Lengkap dengan Middleware dan JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
)

type User struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Address string `json:"address,omitempty"`
}

type CreateUserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func LogMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("From middleware:", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)
    })
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "cuma login yang bisa lanjut")
            return
        }

        next.ServeHTTP(w, r)
    })
}

func StartChiServer(address string) {
    handler := chi.NewRouter()

    handler.Use(LogMiddleware)

    handler.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Hello world")
    })

    handler.Get("/user", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        u := User{
            Name:    "John",
            Age:     17,
            Address: "BSD",
        }

        json.NewEncoder(w).Encode(u)
    })

    handler.Post("/users", func(w http.ResponseWriter, r *http.Request) {
        var payload CreateUserRequest

        err := json.NewDecoder(r.Body).Decode(&payload)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "unable to decode payload")
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)

        json.NewEncoder(w).Encode(payload)
    })

    handler.Group(func(r chi.Router) {
        r.Use(AuthorizationMiddleware)

        r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "This is protected")
        })
    })

    s := http.Server{
        Addr:    address,
        Handler: handler,
    }

    err := s.ListenAndServe()
    if err != nil {
        log.Println(err)
    }
}
```

---

# Intinya

Untuk JSON:

- `json.Marshal` mengubah struct Go menjadi JSON.
- `json.Unmarshal` mengubah JSON menjadi struct Go.
- `json.NewEncoder(w).Encode(data)` cocok untuk menulis JSON response.
- `json.NewDecoder(r.Body).Decode(&data)` cocok untuk membaca JSON request body.
- Field struct harus huruf besar supaya bisa dibaca oleh JSON.
- `json:"-"` membuat field tidak muncul di JSON.
- `omitempty` membuat field kosong tidak muncul di JSON.

Untuk HTTP server:

- `http.ResponseWriter` dipakai untuk membuat response.
- `*http.Request` berisi data request dari client.
- `http.NewServeMux()` adalah router bawaan Go.
- `http.Server` digunakan untuk menjalankan server.

Untuk middleware:

- Middleware adalah wrapper sebelum atau sesudah handler utama.
- Middleware bisa menghentikan request atau melanjutkan ke handler berikutnya.
- Di Chi, lanjut ke handler berikutnya dengan:

```go
next.ServeHTTP(w, r)
```

Untuk framework:

- Gin lebih banyak menyediakan helper lewat `gin.Context`.
- Chi lebih dekat dengan standard `net/http`.
- Keduanya bisa dipasang sebagai `Handler` di `http.Server`.

api.go

func Router(user *User) http.Handler {
    // build construct User

    handler := chi.NewMux()
    handler.Post("/users", user.HandleInsert)
    handler.Get("/users". user.HandleGet) and stuff
}

# Context

`context.Context` digunakan untuk membawa informasi dan sinyal selama proses berjalan.

Dalam HTTP request, context biasanya hidup dari saat request masuk sampai response selesai dikirim.

Context bisa dipakai untuk:

- membawa value kecil yang request-scoped, misalnya `user_id`
- membatalkan proses dengan cancel
- memberi timeout atau deadline
- mendeteksi kalau client sudah disconnect
- meneruskan cancel signal ke database query atau HTTP call lain

---

## Context dengan Value

Contoh sederhana:

```go
ctx := context.Background()

ctx = context.WithValue(ctx, "user_id", "123")

val := ctx.Value("user_id")

s, ok := val.(string)
if !ok {
    fmt.Println("user_id is not string")
    return
}

fmt.Println(s)
```

Penjelasan:

```go
context.WithValue(ctx, "user_id", "123")
```

Artinya kita membuat context baru yang membawa value `user_id`.

```go
ctx.Value("user_id")
```

Digunakan untuk mengambil value dari context.

```go
s, ok := val.(string)
```

Ini adalah type assertion yang aman.

Kalau value ternyata bukan string, `ok` akan bernilai `false`.

---

## Jangan Pakai String sebagai Key Context

Kalau membuat package, jangan pakai string biasa sebagai key context.

Contoh yang kurang aman:

```go
ctx = context.WithValue(ctx, "user_id", "123")
```

Kenapa kurang aman?

Karena package lain juga bisa memakai key `"user_id"`, sehingga value-nya bisa ketimpa.

Lebih aman pakai custom type yang private di package tersebut.

```go
type ctxKeyUserID struct{}
```

Contoh lengkap:

```go
package authctx

import "context"

type ctxKeyUserID struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
    val := ctx.Value(ctxKeyUserID{})

    userID, ok := val.(string)
    return userID, ok
}
```

Pemakaian:

```go
ctx := context.Background()

ctx = authctx.WithUserID(ctx, "123")

userID, ok := authctx.UserIDFromContext(ctx)
if !ok {
    fmt.Println("user_id not found")
    return
}

fmt.Println(userID)
```

Catatan:

`ctxKeyUserID` dibuat private karena huruf awalnya kecil.

Jadi package lain tidak bisa asal membuat key yang sama.

---

## WithCancel

`context.WithCancel` digunakan untuk membuat context yang bisa dibatalkan manual.

Contoh:

```go
ctx := context.Background()

ctx, cancel := context.WithCancel(ctx)

go func() {
    time.Sleep(3 * time.Second)
    cancel()
}()

fmt.Println("Waiting for cancellation")

<-ctx.Done()

fmt.Println(ctx.Err())
```

Penjelasan:

```go
ctx, cancel := context.WithCancel(ctx)
```

Membuat context baru yang bisa dibatalkan.

```go
cancel()
```

Membatalkan context.

```go
<-ctx.Done()
```

Menunggu sampai context dibatalkan.

```go
ctx.Err()
```

Mengambil alasan kenapa context selesai.

Biasanya hasilnya:

```text
context canceled
```

---

## WithTimeout

`context.WithTimeout` digunakan untuk membuat context yang otomatis cancel setelah durasi tertentu.

```go
ctx := context.Background()

ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()

fmt.Println("Waiting for timeout")

<-ctx.Done()

fmt.Println(ctx.Err())
```

Kalau timeout habis, output `ctx.Err()` biasanya:

```text
context deadline exceeded
```

---

## WithDeadline

`context.WithDeadline` mirip dengan `WithTimeout`, tapi dia memakai waktu spesifik.

```go
ctx := context.Background()

deadline := time.Now().Add(3 * time.Second)

ctx, cancel := context.WithDeadline(ctx, deadline)
defer cancel()

<-ctx.Done()

fmt.Println(ctx.Err())
```

Bedanya:

```go
context.WithTimeout(ctx, 3*time.Second)
```

Artinya cancel setelah 3 detik.

```go
context.WithDeadline(ctx, time.Now().Add(3*time.Second))
```

Artinya cancel pada waktu tertentu.

---

## select dengan ctx.Done

Kalau ada proses blocking, kita bisa pakai `select` untuk melihat mana yang selesai duluan.

Contoh:

```go
select {
case <-ctx.Done():
    fmt.Println("context canceled:", ctx.Err())

case val := <-ch:
    fmt.Println("received:", val)
}
```

Artinya:

- kalau context selesai duluan, jalankan case `ctx.Done()`
- kalau channel `ch` menerima data duluan, jalankan case `val := <-ch`

Ini berguna supaya goroutine tidak menunggu selamanya.

---

## Context di HTTP Request

Dalam HTTP handler, request sudah punya context bawaan.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    select {
    case <-ctx.Done():
        fmt.Println("request canceled:", ctx.Err())
        return

    default:
        fmt.Println("request still active")
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

Context dari `r.Context()` bisa selesai kalau:

- request selesai
- client disconnect
- server cancel request
- timeout terjadi

Jadi kalau user menekan `Ctrl+C` di client atau koneksi putus, server bisa mendeteksi lewat:

```go
<-ctx.Done()
```

---

## Context untuk Database

Context sering dikirim ke database query.

```go
func GetUser(ctx context.Context, db *sql.DB, id int) error {
    row := db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = ?", id)

    var userID int

    err := row.Scan(&userID)
    if err != nil {
        return err
    }

    return nil
}
```

Kalau request client sudah cancel, context dari `r.Context()` juga bisa ikut membatalkan query database.

Contoh di HTTP handler:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    err := GetUser(ctx, db, 1)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
```

---

## context.WithoutCancel

`context.WithoutCancel` membuat context baru yang tidak ikut cancel walaupun parent context-nya cancel.

Contoh:

```go
ctx := context.Background()

ctx, cancel := context.WithCancel(ctx)

newCtx := context.WithoutCancel(ctx)

cancel()

fmt.Println(ctx.Err())
fmt.Println(newCtx.Err())
```

Context lama akan canceled, tapi `newCtx` tidak ikut canceled.

Catatan penting:

Kalau memakai `WithoutCancel`, `Done()` bisa bernilai `nil`.

Jadi kalau kita menulis:

```go
<-newCtx.Done()
```

program bisa menunggu selamanya.

Karena itu hati-hati pakai `WithoutCancel`.

Biasanya dipakai kalau ada proses background yang tetap harus lanjut walaupun request utama sudah selesai.

---

## Gin Context vs context.Context

Di Gin, handler menerima:

```go
func(ctx *gin.Context)
```

`gin.Context` bukan `context.Context`.

`gin.Context` adalah context milik Gin yang berisi helper untuk request dan response.

Contoh:

```go
ctx.String(http.StatusOK, "Hello world")
ctx.JSON(http.StatusOK, data)
ctx.Request
ctx.Writer
```

Kalau butuh context standar Go, ambil dari request:

```go
goCtx := ctx.Request.Context()
```

Contoh untuk database:

```go
handler.GET("/users", func(ctx *gin.Context) {
    goCtx := ctx.Request.Context()

    err := db.QueryRowContext(goCtx, "SELECT ...").Scan(...)
    if err != nil {
        ctx.String(http.StatusInternalServerError, err.Error())
        return
    }

    ctx.String(http.StatusOK, "OK")
})
```

---

# HTTP Client

HTTP client digunakan untuk mengirim request ke server lain.

Di Go, HTTP client bawaan ada di package:

```go
net/http
```

---

## Contoh HTTP Client Sederhana

```go
func ClientExample() {
    cl := http.Client{}

    req, err := http.NewRequest(
        http.MethodGet,
        "http://127.0.0.1:8080/hello",
        nil,
    )
    if err != nil {
        log.Println(err)
        return
    }

    req.Header.Set("Content-Type", "application/json")

    res, err := cl.Do(req)
    if err != nil {
        log.Println(err)
        return
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        log.Println("unexpected status:", res.StatusCode)
        return
    }

    b, err := io.ReadAll(res.Body)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(string(b))
}
```

Penjelasan:

```go
cl := http.Client{}
```

Membuat HTTP client.

```go
http.NewRequest(...)
```

Membuat request.

```go
req.Header.Set(...)
```

Menambahkan header.

```go
cl.Do(req)
```

Mengirim request ke server dan menerima response.

```go
defer res.Body.Close()
```

Menutup response body setelah selesai dibaca.

```go
io.ReadAll(res.Body)
```

Membaca isi response body.

---

## HTTP Client dengan Context

Lebih bagus pakai context supaya request bisa cancel atau timeout.

```go
func ClientWithContext() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodGet,
        "http://127.0.0.1:8080/hello",
        nil,
    )
    if err != nil {
        log.Println(err)
        return
    }

    cl := http.Client{}

    res, err := cl.Do(req)
    if err != nil {
        log.Println(err)
        return
    }
    defer res.Body.Close()

    b, err := io.ReadAll(res.Body)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(string(b))
}
```

Kalau request lebih dari 5 detik, context akan cancel dan `cl.Do(req)` akan return error.

---

## POST JSON dengan HTTP Client

Kalau ingin mengirim JSON ke server, gunakan `json.Marshal` lalu masukkan ke body request.

```go
type UserRequest struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func PostJSONExample() {
    ctx := context.Background()

    payload := UserRequest{
        Name: "John",
        Age:  17,
    }

    b, err := json.Marshal(payload)
    if err != nil {
        log.Println(err)
        return
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        "http://127.0.0.1:8080/users",
        bytes.NewReader(b),
    )
    if err != nil {
        log.Println(err)
        return
    }

    req.Header.Set("Content-Type", "application/json")

    cl := http.Client{}

    res, err := cl.Do(req)
    if err != nil {
        log.Println(err)
        return
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusCreated {
        log.Println("unexpected status:", res.StatusCode)
        return
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(string(body))
}
```

---

## Membaca Response JSON

Kalau response dari server berupa JSON, gunakan `json.NewDecoder`.

```go
type UserResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func ReadJSONResponse() {
    ctx := context.Background()

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodGet,
        "http://127.0.0.1:8080/users/1",
        nil,
    )
    if err != nil {
        log.Println(err)
        return
    }

    cl := http.Client{}

    res, err := cl.Do(req)
    if err != nil {
        log.Println(err)
        return
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        log.Println("unexpected status:", res.StatusCode)
        return
    }

    var user UserResponse

    err = json.NewDecoder(res.Body).Decode(&user)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println(user.ID, user.Name, user.Age)
}
```

Catatan:

Untuk membaca request atau response body JSON:

```go
json.NewDecoder(res.Body).Decode(&user)
```

Bukan:

```go
json.NewEncoder(res.Body)
```

`Encoder` dipakai untuk menulis JSON.

`Decoder` dipakai untuk membaca JSON.

---

# Wrapper HTTP Client

Kalau sering melakukan HTTP request, lebih rapi kalau kita membuat wrapper client sendiri.

Contoh struktur file:

```text
project/
├── main.go
├── httpclient/
│   └── client.go
└── users/
    └── client.go
```

---

## httpclient/client.go

File ini berisi wrapper dasar untuk HTTP request.

```go
package httpclient

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type Client struct {
    baseURL string
    client  *http.Client
}

func New(baseURL string, client *http.Client) *Client {
    if client == nil {
        client = http.DefaultClient
    }

    return &Client{
        baseURL: baseURL,
        client:  client,
    }
}

func (c *Client) JSON(
    ctx context.Context,
    method string,
    path string,
    payload any,
    out any,
) (*http.Response, error) {
    var body io.Reader

    if payload != nil {
        b, err := json.Marshal(payload)
        if err != nil {
            return nil, err
        }

        body = bytes.NewReader(b)
    }

    req, err := http.NewRequestWithContext(
        ctx,
        method,
        c.baseURL+path,
        body,
    )
    if err != nil {
        return nil, err
    }

    req.Header.Set("Accept", "application/json")

    if payload != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    res, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if out != nil {
        err = json.NewDecoder(res.Body).Decode(out)
        if err != nil && err != io.EOF {
            return res, err
        }
    }

    return res, nil
}
```

Penjelasan:

```go
type Client struct {
    baseURL string
    client  *http.Client
}
```

Struct ini menyimpan base URL dan HTTP client asli.

```go
func New(baseURL string, client *http.Client) *Client
```

Constructor untuk membuat wrapper client.

```go
func (c *Client) JSON(...)
```

Method reusable untuk mengirim HTTP request dengan JSON.

Bagian ini:

```go
json.Marshal(payload)
```

Mengubah payload Go menjadi JSON.

Bagian ini:

```go
http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
```

Membuat HTTP request dengan context.

Bagian ini:

```go
c.client.Do(req)
```

Mengirim request ke server.

Bagian ini:

```go
json.NewDecoder(res.Body).Decode(out)
```

Membaca response JSON ke variable `out`.

---

# Users Client

`users/client.go` adalah client khusus untuk fitur user.

Dia memakai wrapper `httpclient.Client`.

---

## users/client.go

```go
package users

import (
    "context"
    "errors"
    "net/http"

    "your-module/httpclient"
)

type Client struct {
    baseClient *httpclient.Client
}

func New(baseClient *httpclient.Client) *Client {
    return &Client{
        baseClient: baseClient,
    }
}

type InsertRequest struct {
    ID   string `json:"id,omitempty"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type InsertResponse struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Message string `json:"message,omitempty"`
}

func (cl *Client) Insert(ctx context.Context, name string, age int) (InsertResponse, error) {
    payload := InsertRequest{
        Name: name,
        Age:  age,
    }

    var response InsertResponse

    res, err := cl.baseClient.JSON(
        ctx,
        http.MethodPost,
        "/users/insert",
        payload,
        &response,
    )
    if err != nil {
        return InsertResponse{}, err
    }

    if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
        return InsertResponse{}, errors.New("error response")
    }

    return response, nil
}
```

Penjelasan:

```go
type Client struct {
    baseClient *httpclient.Client
}
```

Client ini adalah wrapper khusus untuk endpoint users.

```go
func New(baseClient *httpclient.Client) *Client
```

Constructor untuk membuat users client.

```go
type InsertRequest struct
```

Struct untuk request body.

```go
type InsertResponse struct
```

Struct untuk response body.

```go
func (cl *Client) Insert(...)
```

Method untuk memanggil endpoint `/users/insert`.

---

# Contoh Cobra Command untuk HTTP Client

Misalnya kita ingin menambahkan command:

```bash
go run main.go client
```

Contoh:

```go
func init() {
    rootCmd.AddCommand(&cobra.Command{
        Use:   "client",
        Short: "Run HTTP client example",
        Run: func(cmd *cobra.Command, args []string) {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()

            base := httpclient.New(
                "http://127.0.0.1:8080",
                http.DefaultClient,
            )

            userClient := users.New(base)

            res, err := userClient.Insert(ctx, "John", 17)
            if err != nil {
                log.Println(err)
                return
            }

            fmt.Println(res)
        },
    })
}
```

---

# Contoh Server untuk Ditest Client

Supaya client di atas bisa dites, server perlu punya route `/users/insert`.

Contoh dengan Chi:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
)

type InsertRequest struct {
    ID   string `json:"id,omitempty"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type InsertResponse struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Message string `json:"message,omitempty"`
}

func main() {
    r := chi.NewRouter()

    r.Post("/users/insert", func(w http.ResponseWriter, r *http.Request) {
        var payload InsertRequest

        err := json.NewDecoder(r.Body).Decode(&payload)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "unable to decode payload")
            return
        }

        response := InsertResponse{
            ID:      "1",
            Name:    payload.Name,
            Age:     payload.Age,
            Message: "user inserted",
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)

        json.NewEncoder(w).Encode(response)
    })

    err := http.ListenAndServe("127.0.0.1:8080", r)
    if err != nil {
        log.Println(err)
    }
}
```

---

# Intinya

Untuk context:

- `context.Context` bisa membawa value dan cancel signal.
- Jangan pakai string biasa sebagai key context di package.
- Pakai custom private type sebagai context key.
- `ctx.Done()` dipakai untuk mendeteksi context selesai atau cancel.
- `context.WithCancel` membuat context yang bisa dibatalkan manual.
- `context.WithTimeout` membuat context yang otomatis cancel setelah durasi tertentu.
- `context.WithDeadline` membuat context yang otomatis cancel pada waktu tertentu.
- `context.WithoutCancel` membuat context baru yang tidak ikut cancel parent, tapi hati-hati karena `Done()` bisa nil.
- Di HTTP handler, pakai `r.Context()` untuk mengambil context request.
- Di Gin, ambil context standar Go dari `ctx.Request.Context()`.

Untuk HTTP client:

- `http.Client{}` digunakan untuk membuat client.
- `http.NewRequestWithContext()` membuat request dengan context.
- `cl.Do(req)` mengirim request ke server dan menerima response.
- `defer res.Body.Close()` wajib supaya body ditutup setelah selesai dibaca.
- `json.NewEncoder(w).Encode(data)` untuk menulis JSON.
- `json.NewDecoder(res.Body).Decode(&data)` untuk membaca JSON.
- Wrapper client membuat kode HTTP request lebih rapi dan reusable.