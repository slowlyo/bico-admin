/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(
  initialState: { currentUser?: API.CurrentUser } | undefined,
) {
  const { currentUser } = initialState ?? {};
  const permissions = currentUser?.permissions || [];
  
  const accessObj: Record<string, boolean> = {};
  
  permissions.forEach((permission: string) => {
    accessObj[permission] = true;
  });
  
  return accessObj;
}
