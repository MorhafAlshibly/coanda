import { ZodIssueCode } from "zod";

export const IssueCode = {
  ...ZodIssueCode,
  not_found: "not_found",
  syntax_error: "syntax_error",
};

export const IssueStatus = {
  invalid_type: 400,
  invalid_literal: 400,
  custom: 400,
  invalid_union: 400,
  invalid_union_discriminator: 400,
  invalid_enum_value: 400,
  unrecognized_keys: 400,
  invalid_arguments: 400,
  invalid_return_type: 400,
  invalid_date: 400,
  invalid_string: 400,
  too_small: 400,
  too_big: 400,
  invalid_intersection_types: 400,
  not_multiple_of: 400,
  not_finite: 400,
  syntax_error: 400,
  not_found: 404,
};

export type IssueCode = keyof typeof IssueCode;
