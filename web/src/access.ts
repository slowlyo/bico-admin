/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(
  initialState: { currentUser?: API.CurrentUser; appConfig?: API.AppConfig } | undefined,
) {
  const { currentUser } = initialState ?? {};
  const permissions = currentUser?.permissions || [];
  
  const accessObj: Record<string, boolean> = {};
  
  permissions.forEach((permission: string) => {
    accessObj[permission] = true;
  });
  accessObj.developer = !!initialState?.appConfig?.debug;
  
  return accessObj;
}
