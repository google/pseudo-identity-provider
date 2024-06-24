/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import {FormlyFieldConfig} from '@ngx-formly/core';

// findAndInsertExpressions updates Formly fields to contain the JSON Schema specified expressions.
export function findAndInsertExpressions(
  field: FormlyFieldConfig,
  schema: any,
) {
  let exp = findExpressions(schema, '');
  insertExpressions(exp, field, '');
}

// findAndInsertExpressionsForArray updates Formly fields to contain the JSON Schema specified
// expressions for array fields.
export function findAndInsertExpressionsForArray(
  field: FormlyFieldConfig,
  schema: any,
) {
  let exp = findExpressions(schema, '');
  let parent = field.parent;

  while (parent != null && parent.key != null) {
    parent = parent.parent;
  }
  insertExpressions(exp, parent!, '');
}

// findExpressions finds fields with expressions and creates a mapping of expressions
// to their fields.
export function findExpressions(
  schema: any,
  path: string,
): Map<string, Map<string, string>> {
  let expressions = new Map<string, Map<string, string>>();
  if (schema == null) {
    return expressions;
  }

  for (const key in schema) {
    // Currently only the hide expression is supported.
    if (key === 'hide') {
      const expression = new Map<string, string>();
      expression.set(key, schema[key]);
      expressions.set(path, expression);
    }
  }

  for (const propKey in schema.properties) {
    const subPath = path === '' ? propKey : path + '.' + propKey;
    const subExpressions = findExpressions(schema.properties[propKey], subPath);
    expressions = new Map([
      ...expressions.entries(),
      ...subExpressions.entries(),
    ]);
  }

  if (schema.items != null) {
    const subExpressions = findExpressions(schema.items, path);
    expressions = new Map([
      ...expressions.entries(),
      ...subExpressions.entries(),
    ]);
  }

  return expressions;
}

// insertExpressions converts the JSON Schema expressions to Formly function expressions
// and inserts them into the FormlyField.
function insertExpressions(
  expressions: Map<string, Map<string, string>>,
  fieldConfig: FormlyFieldConfig,
  path: string,
) {
  if (fieldConfig.fieldGroup == null) {
    return;
  }

  for (const field of fieldConfig.fieldGroup) {
    let fieldPath = path === '' ? String(field.key) : path + '.' + field.key;

    // Array object use indices for keys, for those we want to use the full field
    // path for matching against the expression map, but use the parent field
    // path to pass on, as that is what is represented in the schema, not the indices.
    let parentIsArray = field.parent != null && field.parent.type === 'array';
    let subPath = parentIsArray ? path : fieldPath;

    const fieldExpressions = expressions.get(fieldPath);
    if (fieldExpressions != null) {
      field.expressions = {};
      for (const [exprKey, expr] of fieldExpressions) {
        if (expr != null) {
          field.expressions[exprKey] = (field: any) => {
            return expressionEval(expr, field);
          };
        }
      }
    }

    insertExpressions(expressions, field, subPath);
  }
}

// expressionEval takes one JSON schema expression and evaluates
// the expression in a manner compatible with Formly. Formly also
// supports its own expression evaluation, but that relies on eval()
// and is not compatible with CSP so we do our own, much more limited,
// evaluation here.
// Expressions can be of the form field_path === value or field_path !== value.
function expressionEval(expr: string, field: FormlyFieldConfig): boolean {
  const equals = expr.includes('===');
  const exprArray = expr.split(equals ? '===' : '!==');
  if (exprArray.length !== 2) {
    throw Error('Invalid Expression: ' + expr);
  }

  const varArray = exprArray[0].trim().split('.');
  const value = exprArray[1].trim();

  let i = 0;
  for (; varArray[i] === 'parent' && i < varArray.length; i++) {
    if (field.parent) {
      field = field.parent;
    }
  }

  let varValue = field.parent?.model;
  for (; i < varArray.length; i++) {
    varValue = varValue[varArray[i].trim()];
  }
  return equals === (varValue === value);
}
