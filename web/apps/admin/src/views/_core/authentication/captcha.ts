import type { InjectionKey } from 'vue';

export interface CaptchaState {
  id: string;
  image: string;
}

export interface CaptchaContext {
  loadCaptcha: () => Promise<void>;
  state: CaptchaState;
}

export const CAPTCHA_KEY: InjectionKey<CaptchaContext> = Symbol('login-captcha');
