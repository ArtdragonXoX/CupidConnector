<template>
  <el-empty v-if="isNoData" style="max-height: 320px">
    <ElButton @click="handleClick" style="margin-top: 10px;" :icon="Refresh" :loading="isLoading">刷新</ElButton>
  </el-empty>
  <div v-if="!isNoData">
    <el-row>
      <el-col :span="8">
        <el-statistic :value="info?.expiration_time">
          <template #title>
            <div style="display: inline-flex; align-items: center">
              <el-icon style="margin-right: 4px" :size="12">
                <Calendar />
              </el-icon>
              过期时间
            </div>
          </template>
        </el-statistic>
      </el-col>
      <el-col :span="8">
        <el-statistic :value="info?.user_name">
          <template #title>
            <div style="display: inline-flex; align-items: center">
              <el-icon style="margin-right: 4px" :size="12">
                <User />
              </el-icon>
              用户名称
            </div>
          </template>
        </el-statistic>
        <ElButton @click="handleClick" style="margin-top: 10px;" :icon="Refresh" :loading="isLoading" size="small">刷新
        </ElButton>
      </el-col>
      <el-col :span="8">
        <el-statistic :value="info?.account_status">
          <template #title>
            <div style="display: inline-flex; align-items: center">
              <el-icon style="margin-right: 4px" :size="12">
                <CircleCheck />
              </el-icon>
              用户状态
            </div>
          </template>
        </el-statistic>
      </el-col>
    </el-row>
    <div class="flux-chart">
      <div class="statistics-area-chart" ref="chartDom"></div>
    </div>
    <div class="surplus" v-if="info">
      剩余流量:{{ parseFloat((info.overall - info?.used).toFixed(3)) }}MB
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElCol, ElRow } from 'element-plus'
import { GetInfo, RefreshInfo } from '../../wailsjs/go/main/App'
import { model } from '../../wailsjs/go/models'
import * as echarts from 'echarts/core';
import { GaugeChart, GaugeSeriesOption } from 'echarts/charts';
import { CanvasRenderer } from 'echarts/renderers';
import { Refresh, Calendar, User, CircleCheck } from "@element-plus/icons-vue";
import { EventsOn } from "../../wailsjs/runtime";
echarts.use([GaugeChart, CanvasRenderer]);
type EChartsOption = echarts.ComposeOption<GaugeSeriesOption>;

const isLoading = ref(false)
const info = ref<model.Info>()
const isNoData = ref<boolean>(true)

const getInfo = async () => {
  isLoading.value = true
  await GetInfo().then((res) => {
    if (res.code === 1) {
      info.value = res.data
      isNoData.value = false;
    } else {
      info.value = {
        user_name: "未知用户",
        overall: 0,
        used: 0,
        expiration_time: "未知",
        account_status: "未知",
      };
      isNoData.value = true;
    }
    if (info.value == null || info.value?.user_name == "") {
      info.value = {
        user_name: "未知用户",
        overall: 0,
        used: 0,
        expiration_time: "未知",
        account_status: "未知",
      };
      isNoData.value = true;
    }
    isLoading.value = false
  })
}

const chartDom = ref<HTMLElement | null>(null);
const myChart = ref<echarts.ECharts | null>(null);
let option: EChartsOption;

const updateOption = async () => {
  if (myChart.value) {
    option = {
      series: [
        {
          min: 0,
          max: info.value?.overall,
          splitNumber: 8,
          type: 'gauge',
          itemStyle: {
            color: '#58D9F9',
            shadowColor: 'rgba(0,138,255,0.45)',
            shadowBlur: 10,
            shadowOffsetX: 2,
            shadowOffsetY: 2,
          },
          progress: {
            show: true,
            width: 18
          },
          axisLine: {
            lineStyle: {
              width: 18,
              color: [
                [0.25, '#7CFFB2'],
                [0.5, '#58D9F9'],
                [0.75, '#FDDD60'],
                [1, '#FF6E76']
              ]
            }
          },
          axisTick: {
            show: false
          },
          splitLine: {
            length: 10,
            lineStyle: {
              width: 2,
              color: '#fff'
            }
          },
          axisLabel: {
            distance: 25,
            color: '#fff',
            fontSize: 10
          },
          anchor: {
            show: true,
            showAbove: true,
            size: 15,
            itemStyle: {
              borderWidth: 10
            }
          },
          title: {
            offsetCenter: [0, '60%'],
          },
          detail: {
            valueAnimation: true,
            fontSize: 20,
            offsetCenter: [0, '85%'],
            formatter: '{value}' + 'MB',
            color: '#fff'
          },
          data: [
            {
              value: info.value?.used,
              name: '已使用',
            }
          ]
        }
      ]
    };
    option && myChart.value.setOption(option)
  }
}

const handleClick = async () => {
  isLoading.value = true
  await RefreshInfo()
  isLoading.value = false
}

const updateInfo = async () => {
  isLoading.value = true
  await getInfo()
  if (chartDom.value) {
    // 初始化 ECharts 实例
    myChart.value = echarts.init(chartDom.value);
    await updateOption();
  }
}

onMounted(async () => {
  EventsOn("updateInfo", async () => {
    isLoading.value = true
    await updateInfo();
  })
})
</script>

<style scoped>
.statistics-area-chart {
  width: 240px;
  height: 240px;
  text-align: center;
}

.flux-chart {
  display: flex;
  justify-content: center;
  text-align: center;
}
</style>