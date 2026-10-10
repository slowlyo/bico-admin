<script lang="ts" setup>
import { inject } from 'vue';

import { CAPTCHA_KEY } from './captcha';

defineOptions({ name: 'CaptchaInput' });

const captcha = inject(CAPTCHA_KEY);
const code = defineModel<string>({ default: '' });

// 点击图片换一张。校验失败后的刷新由登录页负责，这里只处理手动刷新。
function refresh() {
  void captcha?.loadCaptcha();
}
</script>

<template>
  <div class="flex w-full items-center gap-2">
    <input
      v-model="code"
      autocomplete="off"
      class="border-input bg-background text-foreground placeholder:text-muted-foreground/50 h-10 min-w-0 flex-1 rounded-md border px-3 text-sm outline-none"
      placeholder="请输入验证码"
    />
    <button
      type="button"
      class="border-input flex h-10 w-32 shrink-0 cursor-pointer items-center justify-center overflow-hidden rounded-md border bg-zinc-100 p-0.5 dark:bg-zinc-400"
      @click="refresh"
    >
      <img
        v-if="captcha?.state.image"
        :src="captcha.state.image"
        alt="验证码"
        class="h-full w-full object-contain"
      />
    </button>
  </div>
</template>
