var editor = CodeMirror.fromTextArea(document.getElementById("code"), {
  lineNumbers: true,     // 显示行数
  indentUnit: 2,         // 缩进单位为2
  styleActiveLine: true, // 当前行背景高亮
  matchBrackets: true,   // 括号匹配
  mode: 'yaml',     // HMTL混合模式
  tabSize: 2,
  lineWrapping: true,    // 自动换行
  theme: 'default',      // 使用monokai模版
  showCursorWhenSelecting : true,
});
editor.setOption("extraKeys", {
  // Tab键换成2个空格
  Tab: function(cm) {
      var spaces = Array(cm.getOption("indentUnit") + 1).join(" ");
      cm.replaceSelection(spaces);
  },
  // F11键切换全屏
  "F11": function(cm) {
      cm.setOption("fullScreen", !cm.getOption("fullScreen"));
  },
  // Esc键退出全屏
  "Esc": function(cm) {
      if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
  }
});
