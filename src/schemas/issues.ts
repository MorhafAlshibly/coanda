/**
 * @openapi
 * components:
 *  schemas:
 *    IssueStatus:
 *      type: number
 *      oneOf:
 *        - syntax_error
 *          const: 400
 *          description: Bad syntax
 *        - not_found
 *          const: 404
 *          description: Resource not found
 */
export enum IssueStatus {
  not_found = 404,
  syntax_error = 400,
}

export type IssueCode = keyof typeof IssueStatus;
