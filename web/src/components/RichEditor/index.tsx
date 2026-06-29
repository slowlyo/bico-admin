import '@wangeditor/editor/dist/css/style.css';

import { Editor, Toolbar } from '@wangeditor/editor-for-react';
import type { IEditorConfig, IToolbarConfig } from '@wangeditor/editor';
import React, { useEffect, useMemo, useState } from 'react';

import { uploadForEditor } from '@/services/common/upload';

export interface RichEditorProps {
  /** HTML 内容 */
  value?: string;
  /** 内容变更回调 */
  onChange?: (html: string) => void;
  /** 是否禁用 */
  disabled?: boolean;
  /** 编辑器高度 */
  height?: number;
}

/**
 * 通用富文本编辑器组件（wangeditor v5）
 */
const RichEditor: React.FC<RichEditorProps> = ({
  value,
  onChange,
  disabled,
  height = 420,
}) => {
  const [editor, setEditor] = useState<any>(null);

  useEffect(() => {
    return () => {
      if (editor == null) return;
      editor.destroy();
      setEditor(null);
    };
  }, [editor]);

  const toolbarConfig = useMemo<Partial<IToolbarConfig>>(() => {
    return {
      excludeKeys: [],
    };
  }, []);

  const editorConfig = useMemo<Partial<IEditorConfig>>(() => {
    return {
      placeholder: '请输入内容...',
      readOnly: !!disabled,
      MENU_CONF: {
        uploadImage: {
          async customUpload(file: File, insertFn: (url: string, alt?: string, href?: string) => void) {
            const resp = await uploadForEditor(file, 'image');
            if (!resp || resp.errno !== 0 || !resp.data?.url) {
              throw new Error(resp?.message || '图片上传失败');
            }
            insertFn(resp.data.url);
          },
        },
        uploadVideo: {
          async customUpload(file: File, insertFn: (url: string) => void) {
            const resp = await uploadForEditor(file, 'video');
            if (!resp || resp.errno !== 0 || !resp.data?.url) {
              throw new Error(resp?.message || '视频上传失败');
            }
            insertFn(resp.data.url);
          },
        },
      },
    };
  }, [disabled]);

  return (
    <div style={{ border: '1px solid #d9d9d9', borderRadius: 6, overflow: 'hidden' }}>
      <Toolbar editor={editor} defaultConfig={toolbarConfig} mode="default" />
      <Editor
        defaultConfig={editorConfig}
        value={value}
        onCreated={setEditor}
        onChange={(ed) => {
          // 这里通过 html 作为外部值，便于和表单/后端存储对接
          const html = ed.getHtml();
          onChange?.(html);
        }}
        mode="default"
        style={{ height, overflowY: 'hidden' }}
      />
    </div>
  );
};

export default RichEditor;
