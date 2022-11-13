import { FC } from 'react';

import { AppBreadcrumb, AppBreadcrumbLink } from '@/components';

type EnglishSentenceMemorizationBreadcrumbProps = {
  breadcrumbLinks: AppBreadcrumbLink[];
  workbookUrl: string;
  name: string;
  id: number;
};

export const EnglishSentenceMemorizationBreadcrumb: FC<
  EnglishSentenceMemorizationBreadcrumbProps
> = (props: EnglishSentenceMemorizationBreadcrumbProps) => {
  const links = [...props.breadcrumbLinks];
  links.push(new AppBreadcrumbLink(props.workbookUrl + props.id, props.name));
  return <AppBreadcrumb links={links} text={'Memorization'} />;
};
