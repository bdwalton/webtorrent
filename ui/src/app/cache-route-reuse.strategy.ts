import { RouteReuseStrategy } from '@angular/router/';
import { ActivatedRouteSnapshot, DetachedRouteHandle } from '@angular/router';

export class CacheRouteReuseStrategy implements RouteReuseStrategy {
  storedRouteHandles = new Map<string, DetachedRouteHandle>();
  allowRetriveCache: Map<string, boolean> = new Map<string, boolean>([
    ['torrent', true],
  ]);

  shouldReuseRoute(before: ActivatedRouteSnapshot, curr:  ActivatedRouteSnapshot): boolean {
    if (this.getPath(curr) === 'torrent') {
      return true;
    }

    return false
  }

  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle | null {
    return this.storedRouteHandles.get(this.getPath(route)) as DetachedRouteHandle;
  }

  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    const path = this.getPath(route);
    if (this.allowRetriveCache.get(path)) {
      return this.storedRouteHandles.has(path);
    }

    return false;
  }

  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    const path = this.getPath(route);
    if (this.allowRetriveCache.get(path)) {
      return true;
    }

    return false;
  }

  store(route: ActivatedRouteSnapshot, detachedTree: DetachedRouteHandle): void {
    this.storedRouteHandles.set(this.getPath(route), detachedTree);
  }

  private getPath(route: ActivatedRouteSnapshot): string {
    if (route.routeConfig !== null && route.routeConfig.path !== null) {
      // non-null assertiong (!) here because we've just validated that.
      return route.routeConfig.path!;
    }

    return '';
  }
}
