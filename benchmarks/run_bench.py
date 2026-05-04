import subprocess
import os
import time
import math
import csv

# Настройки
EXECUTABLE = "./huff"
RUNS = 15
TEST_DIR = "benchmarks/testdata"
RESULT_CSV = "benchmarks/results.csv"

def get_file_size(path):
    return os.path.getsize(path)

def benchmark_file(filepath):
    orig_size = get_file_size(filepath)
    out_huf = filepath + ".huf"
    out_rest = filepath + ".rest"
    
    comp_times = []
    decomp_times = []

    for _ in range(RUNS):
        # Замер сжатия
        start = time.perf_counter()
        subprocess.run([EXECUTABLE, "-c", filepath, out_huf], capture_output=True)
        comp_times.append(time.perf_counter() - start)

        # Замер разжатия
        start = time.perf_counter()
        subprocess.run([EXECUTABLE, "-d", out_huf, out_rest], capture_output=True)
        decomp_times.append(time.perf_counter() - start)

    comp_size = get_file_size(out_huf)
    
    # Удаляем временные файлы
    os.remove(out_huf)
    os.remove(out_rest)

    # Подсчёт матожидания и среднеквадратичного отклонения
    def calc_stats(times):
        mean = sum(times) / len(times)
        variance = sum((x - mean) ** 2 for x in times) / len(times)
        std = math.sqrt(variance)
        return mean, std

    comp_mean, comp_std = calc_stats(comp_times)
    decomp_mean, decomp_std = calc_stats(decomp_times)

    ratio = (comp_size / orig_size) * 100

    orig_mb = orig_size / (1024 * 1024)
    comp_speed = orig_mb / comp_mean if comp_mean > 0 else 0
    decomp_speed = orig_mb / decomp_mean if decomp_mean > 0 else 0

    return {
        "ratio": ratio,
        "comp_speed": comp_speed, "comp_std": (orig_mb / comp_mean - orig_mb / (comp_mean + comp_std)) if comp_mean > 0 else 0, # приближенная ошибка скорости
        "decomp_speed": decomp_speed, "decomp_std": (orig_mb / decomp_mean - orig_mb / (decomp_mean + decomp_std)) if decomp_mean > 0 else 0
    }

def main():
    files = []
    for f in os.listdir(TEST_DIR):
        if "1kb" in f or "10kb" in f or "100kb" in f or "1024kb" in f:
            files.append(os.path.join(TEST_DIR, f))

    files.sort(key = lambda x: os.path.getsize(x))

    with open(RESULT_CSV, "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["File", "Type", "SizeBytes", "RatioPct", "CompSpeedMBs", "CompStdMBs", "DecompSpeedMBs", "DecompStdMBs"])

        for filepath in files:
            print(f"Testing: {filepath}...")
            stats = benchmark_file(filepath)
            
            f_name = os.path.basename(filepath)
            f_type = "text" if "text" in f_name else ("random" if "rand" in f_name else "compressed")
            size = get_file_size(filepath)

            writer.writerow([f_name, f_type, size, f"{stats['ratio']:.2f}", 
                             f"{stats['comp_speed']:.2f}", f"{abs(stats['comp_std']):.2f}",
                             f"{stats['decomp_speed']:.2f}", f"{abs(stats['decomp_std']):.2f}"])
            
    print(f"Done! Results saved in {RESULT_CSV}")

if __name__ == "__main__":
    main()
