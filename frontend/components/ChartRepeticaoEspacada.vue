
<template>
    <div class="card">
        <Chart type="line" :data="chartData" :options="chartOptions" class="h-[30rem]" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import Chart from 'primevue/chart';

onMounted(() => {
    chartData.value = setChartData();
    chartOptions.value = setChartOptions();
});

const k = 1.84
const c = 1.25
const minutes = [1, 20, 60, 60*9, 60*24, 60*24*6, 60*24*31] // 1min, 20min, 1h, 9h, 1d, 6d, 31d

const memory_retention = (v: Array<number>): Array<number> => {
    return v.map(value => 100 * k / ((Math.log10(value) ** c) + k));
}

const retention = memory_retention(minutes)

const chartData = ref();
const chartOptions = ref();

const setChartData = () => {
    const documentStyle = getComputedStyle(document.documentElement);

    return {
        labels: ['1min', '20min', '1h', '9h', '1d', '6d', '31d'],
        datasets: [
            {
                label: 'Retenção',
                data: retention,
                fill: false,
                borderColor: documentStyle.getPropertyValue('--p-cyan-500'),
                tension: 0.4
            }
        ]
    };
};
const setChartOptions = () => {
    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--p-text-color');
    const textColorSecondary = documentStyle.getPropertyValue('--p-text-muted-color');
    const surfaceBorder = documentStyle.getPropertyValue('--p-content-border-color');

    return {
        maintainAspectRatio: false,
        aspectRatio: 0.6,
        plugins: {
            legend: {
                labels: {
                    color: textColor
                }
            }
        },
        scales: {
            x: {
                ticks: {
                    color: textColorSecondary
                },
                grid: {
                    color: surfaceBorder
                }
            },
            y: {
                ticks: {
                    color: textColorSecondary
                },
                grid: {
                    color: surfaceBorder
                }
            }
        }
    };
}
</script>
