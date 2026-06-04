import {defineRouting} from 'next-intl/routing';
import {createNavigation} from 'next-intl/navigation';
import {UI_LOCALES} from '@/lib/language';
 
export const routing = defineRouting({
  locales: [...UI_LOCALES],
  defaultLocale: 'en'
});
 
// Lightweight wrappers around Next.js' navigation APIs
// that will consider the routing configuration
export const {Link, redirect, usePathname, useRouter, getPathname} =
  createNavigation(routing);
