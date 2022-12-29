import { ZodIssueCode } from "zod";

export const IssueCode = {
  ...ZodIssueCode,
  not_found: "not_found",
  syntax_error: "syntax_error",
};
export type IssueCode = keyof typeof IssueCode;

export const failModel = (code: IssueCode, path: Array<string>, message: string) => {
  return [
    {
      code,
      path,
      message,
    },
  ];
};
