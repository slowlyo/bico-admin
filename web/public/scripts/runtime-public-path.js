// setRuntimePublicPath 从当前脚本地址推导后台静态资源根路径。
(function setRuntimePublicPath() {
  // 脚本由后台入口目录加载，删除自身文件名后即为当前静态资源根路径。
  const script = document.currentScript;
  const source = script && script.src;
  const suffix = '/scripts/runtime-public-path.js';

  if (source && source.includes(suffix)) {
    window.publicPath = source.slice(0, source.indexOf(suffix) + 1);
  }
})();
