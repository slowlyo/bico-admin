<template>
  <div>
    <SectionTitle :title="$t('setting.color.title')" class="mt-10" />
    <div class="-mr-4">
      <div class="flex flex-wrap items-center">
        <div
          v-for="color in configOptions.mainColors"
          :key="color"
          class="flex items-center justify-center size-[23px] mr-4 mb-2.5 cursor-pointer rounded-full transition-all duration-200 hover:opacity-85"
          :style="{ background: `${color} !important` }"
          @click="colorHandlers.selectColor(color)"
        >
          <ArtSvgIcon
            icon="ri:check-fill"
            class="text-base !text-white"
            v-show="color === systemThemeColor"
          />
        </div>

        <!-- 自定义颜色选择 -->
        <div
          class="custom-color-picker flex items-center mb-2.5"
          :class="{ 'is-active': !configOptions.mainColors.includes(systemThemeColor) }"
        >
          <el-color-picker
            v-model="systemThemeColor"
            :predefine="[...configOptions.mainColors]"
            size="small"
            @change="(val: string | null) => val && colorHandlers.selectColor(val)"
          />
          <ArtSvgIcon
            v-if="!configOptions.mainColors.includes(systemThemeColor)"
            icon="ri:check-fill"
            class="check-icon text-base !text-white"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .custom-color-picker {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 23px;
    height: 23px;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: border-box;
    
    // 使用双重背景实现超细五彩边框
    background: 
      linear-gradient(#fff, #fff) padding-box,
      conic-gradient(
        from 0deg,
        #ff0000, #ff7f00, #ffff00, #00ff00, #0000ff, #4b0082, #8b00ff, #ff0000
      ) border-box;
    border: 2px solid transparent;

    &:hover {
      transform: scale(1.08);
    }

    &.is-active {
      transform: scale(1.15);
      box-shadow: 0 0 8px rgba(0, 0, 0, 0.15);
    }

    :deep(.el-color-picker) {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    :deep(.el-color-picker__trigger) {
      width: 100% !important;
      height: 100% !important;
      padding: 0;
      border: none;
      border-radius: 50%;
      background: transparent !important;
    }

    :deep(.el-color-picker__color) {
      width: 100%;
      height: 100%;
      border: none;
      border-radius: 50%;
    }

    :deep(.el-color-picker__color-inner) {
      width: 100%;
      height: 100%;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    :deep(.el-color-picker__icon),
    :deep(.el-color-picker__empty) {
      display: none !important;
    }

    .check-icon {
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      pointer-events: none;
      z-index: 2;
    }
  }
</style>

<script setup lang="ts">
  import SectionTitle from './SectionTitle.vue'
  import { useSettingStore } from '@/store/modules/setting'
  import { useSettingsConfig } from '../composables/useSettingsConfig'
  import { useSettingsHandlers } from '../composables/useSettingsHandlers'
  import { storeToRefs } from 'pinia'

  const settingStore = useSettingStore()
  const { systemThemeColor } = storeToRefs(settingStore)
  const { configOptions } = useSettingsConfig()
  const { colorHandlers } = useSettingsHandlers()
</script>
