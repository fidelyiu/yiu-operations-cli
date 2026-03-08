---
title: TS中类型安全的路径定义
authors: [FidelYiu]
tags: [ts]
---

## TS中类型安全的路径定义

在TS定义类型安全的Path树结构。

主要原理是TS的类型也能支持运算符了。

<!-- truncate -->

## 方式1

### builder 和 类型

`pathBuilder.ts`

```ts
export type PathDefinition =
  | string
  | {
      segment: string;
      children: Record<string, PathDefinition>;
    };

export const group = <
  const TSegment extends string,
  const TChildren extends Record<string, PathDefinition>,
>(
  segment: TSegment,
  children: TChildren,
) => {
  return {
    segment,
    children,
  } as const;
};

type BuildPathResult<
  TBase extends string,
  TDef extends PathDefinition,
> = TDef extends string
  ? `${TBase}/${TDef}`
  : TDef extends {
        segment: infer TSegment extends string;
        children: infer TChildren extends Record<string, PathDefinition>;
      }
    ? {
        root: `${TBase}/${TSegment}`;
      } & {
        [K in keyof TChildren]: BuildPathResult<
          `${TBase}/${TSegment}`,
          TChildren[K]
        >;
      }
    : never;

export const buildPaths = <
  const TBase extends string,
  const TDef extends PathDefinition,
>(
  base: TBase,
  definition: TDef,
): BuildPathResult<TBase, TDef> => {
  if (typeof definition === "string") {
    return `${base}/${definition}` as BuildPathResult<TBase, TDef>;
  }

  const root = `${base}/${definition.segment}`;
  const result: Record<string, unknown> = {
    root,
  };

  Object.keys(definition.children).forEach((key) => {
    result[key] = buildPaths(root, definition.children[key]);
  });

  return result as BuildPathResult<TBase, TDef>;
};
```

### 使用

```ts
const consolePath = buildPaths(
  "",
  group("console", {
    resource: group("resource", {
      dashboard: "dashboard",
      user: group("user", {
        supplierApplication: "supplier-application",
        supplierManagement: "supplier-management",
        supplierBlacklist: "supplier-blacklist",
        demanderApplication: "demander-application",
        demanderManagement: "demander-management",
        demanderBlacklist: "demander-blacklist",
      }),
      mini: group("mini", {
        nonPublicResourceManagement: "non-public-resource-management",
        teacherManagement: "teacher-management",
        publicResourceManagement: "public-resource-management",
        bannerManagement: "banner-management",
        hotResourceConfig: "hot-resource-config",
        interestConfig: "interest-config",
        websiteManagement: "website-management",
      }),
      order: group("order", {
        list: "list",
      }),
      fund: group("fund", {
        overview: "overview",
      }),
      afterSalesDispute: group("after-sales-dispute", {
        list: "after-sales/list",
      }),
      operationRiskControl: group("operation-risk-control", {
        overview: "risk/overview",
      }),
      systemSetting: group("system-setting", {
        basic: "system/basic",
      }),
    }),
  }),
);

export const Path = {
  console: consolePath,
};
```

## 方式2

### builder 和 类型

`pathBuilder.ts`

```ts
export type PathDefinition = {
  segment: string;
  children?: Record<string, PathDefinition>;
};

export const segment = <const TSegment extends string>(segment: TSegment) => {
  return {
    segment,
  } as const;
};

export const group = <
  const TSegment extends string,
  const TChildren extends Record<string, PathDefinition>,
>(
  segment: TSegment,
  children: TChildren,
) => {
  return {
    segment,
    children,
  } as const;
};

type BuildPathResult<
  TBase extends string,
  TDef extends PathDefinition,
> = TDef extends {
  segment: infer TSegment extends string;
  children: infer TChildren extends Record<string, PathDefinition>;
}
  ? {
      root: `${TBase}/${TSegment}`;
    } & {
      [K in keyof TChildren]: BuildPathResult<
        `${TBase}/${TSegment}`,
        TChildren[K]
      >;
    }
  : TDef extends {
        segment: infer TSegment extends string;
      }
    ? `${TBase}/${TSegment}`
    : never;

export const buildPaths = <
  const TBase extends string,
  const TDef extends PathDefinition,
>(
  base: TBase,
  definition: TDef,
): BuildPathResult<TBase, TDef> => {
  const root = `${base}/${definition.segment}`;
  if (!definition.children) {
    return root as BuildPathResult<TBase, TDef>;
  }

  const result: Record<string, unknown> = {
    root,
  };

  Object.keys(definition.children).forEach((key) => {
    result[key] = buildPaths(root, definition.children![key]);
  });

  return result as BuildPathResult<TBase, TDef>;
};
```

### 使用

```ts
import { buildPaths, group, segment } from "@/utils/pathBuilder";

const consolePath = buildPaths(
  "",
  group("console", {
    resource: group("resource", {
      dashboard: segment("dashboard"),
      user: group("user", {
        supplierApplication: segment("supplier-application"),
        supplierManagement: segment("supplier-management"),
        supplierBlacklist: segment("supplier-blacklist"),
        demanderApplication: segment("demander-application"),
        demanderManagement: segment("demander-management"),
        demanderBlacklist: segment("demander-blacklist"),
      }),
      mini: group("mini", {
        nonPublicResourceManagement: segment("non-public-resource-management"),
        teacherManagement: segment("teacher-management"),
        publicResourceManagement: segment("public-resource-management"),
        bannerManagement: segment("banner-management"),
        hotResourceConfig: segment("hot-resource-config"),
        interestConfig: segment("interest-config"),
        websiteManagement: segment("website-management"),
      }),
      order: group("order", {
        list: segment("list"),
      }),
      fund: group("fund", {
        overview: segment("overview"),
      }),
      afterSalesDispute: group("after-sales-dispute", {
        list: segment("after-sales/list"),
      }),
      operationRiskControl: group("operation-risk-control", {
        overview: segment("risk/overview"),
      }),
      systemSetting: group("system-setting", {
        basic: segment("system/basic"),
      }),
    }),
  }),
);

export const Path = {
  console: consolePath,
} as const;
```
