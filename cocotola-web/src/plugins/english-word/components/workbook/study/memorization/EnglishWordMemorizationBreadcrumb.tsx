import React from 'react';

import { AppBreadcrumb, AppBreadcrumbLink } from 'components';

export const EnglishWordMemorizationBreadcrumb: React.FC<EnglishWordMemorizationBreadcrumbProps> = (
  props: EnglishWordMemorizationBreadcrumbProps
) => {
  const links = [...props.breadcrumbLinks];
  links.push(new AppBreadcrumbLink(props.workbookUrl + props.id, props.name,));
  return (
    <AppBreadcrumb
      links={links}
      text={'Memorization'}
    />
  );
};
type EnglishWordMemorizationBreadcrumbProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
  name: string;
  id: number;
};
