<!-- 布局容器 -->
<template>
  <div class="app-layout">
    <aside id="app-sidebar">
      <ArtSidebarMenu />
    </aside>

    <main id="app-main">
      <div id="app-header">
        <ArtHeaderBar />
      </div>
      <div id="app-content">
        <ArtPageContent />
      </div>
    </main>

    <div id="app-global">
      <ArtGlobalComponent />
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'AppLayout' })
</script>

<style lang="scss" scoped>
  @use './style';

  .app-layout {
    #app-main #app-content {
      :deep(.layout-content) {
        box-sizing: border-box;
        width: calc(100% - 40px);
        margin: auto;
      }
    }

    @media only screen and (width <= 800px) {
      #app-main #app-content {
        :deep(.layout-content) {
          width: calc(100% - 40px);
        }
      }
    }

    @media only screen and (width <= 640px) {
      #app-main #app-content {
        :deep(.layout-content) {
          width: calc(100% - 30px);
        }
      }
    }
  }
</style>
