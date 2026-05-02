import pandas as pd
import matplotlib.pyplot as plt

# Чтение данных
df = pd.read_csv("benchmarks/results.csv")

# Настройка стиля графиков
plt.style.use('seaborn-v0_8-whitegrid')

# График 1: Степень сжатия
plt.figure(figsize=(10, 6))
for dtype in df['Type'].unique():
    subset = df[df['Type'] == dtype]
    plt.plot(subset['SizeBytes'], subset['RatioPct'], marker='o', label=dtype, linewidth=2)

# Линия 100% — граница безубыточности
plt.axhline(y=100, color='red', linestyle='--', linewidth=1.5, label='Без сжатия (100%)')

plt.xscale('log') # Логарифмическая шкала по X
plt.xlabel('Размер исходного файла (байт)', fontsize=12)
plt.ylabel('Размер архива (% от исходного)', fontsize=12)
plt.title('Зависимость степени сжатия от размера и типа данных', fontsize=14)
plt.legend(fontsize=11)
plt.tight_layout()
plt.savefig('benchmarks/graphs/graph_ratio.png', dpi=150)
plt.close()

# График 2: Скорость работы
fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 6))

for dtype in df['Type'].unique():
    subset = df[df['Type'] == dtype]
    
    # Сжатие
    ax1.errorbar(subset['SizeBytes'], subset['CompSpeedMBs'], yerr=subset['CompStdMBs'], 
                 marker='o', label=dtype, capsize=5, linewidth=2)
    # Разжатие
    ax2.errorbar(subset['SizeBytes'], subset['DecompSpeedMBs'], yerr=subset['DecompStdMBs'], 
                 marker='o', label=dtype, capsize=5, linewidth=2)

for ax, title in [(ax1, 'Сжатие (Compress)'), (ax2, 'Разжатие (Decompress)')]:
    ax.set_xscale('log')
    ax.set_xlabel('Размер исходного файла (байт)', fontsize=12)
    ax.set_ylabel('Скорость (МБ/с)', fontsize=12)
    ax.set_title(title, fontsize=14)
    ax.legend(fontsize=11)

plt.suptitle('Зависимость скорости работы от размера файла (усы - погрешность измерений, среднеквадратичное отклонение)', fontsize=14, y=0.98)
plt.tight_layout()
plt.savefig('benchmarks/graphs/graph_speed.png', dpi=150)
plt.close()

print("Графики сохранены: benchmarks/graphs/graph_ratio.png и benchmarks/graphs/graph_speed.png")
