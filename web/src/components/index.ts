/**
 * 这个文件作为组件的目录
 * 目的是统一管理对外输出的组件，方便分类
 */
/**
 * 布局组件
 */
import Footer from './Footer';
import { Question, SelectLang } from './RightContent';
import { AvatarDropdown, AvatarName } from './RightContent/AvatarDropdown';
import PageContainer from './PageContainer';

import RichEditor from './RichEditor';

import ExcelImportExportActions from './ExcelImportExportActions';

/**
 * CRUD 组件
 */
import CrudModal from './CrudModal';
import CrudTable from './CrudTable';

export {
  AvatarDropdown,
  AvatarName,
  Footer,
  PageContainer,
  Question,
  SelectLang,
  CrudModal,
  CrudTable,
  RichEditor,
  ExcelImportExportActions,
};
