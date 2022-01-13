variable "role_name" {
  type        = string
  description = "The name to give the role"
}

variable "iam_path" {
  type        = string
  default     = "/service/"
  description = "The path to give to the IAM role (default: /service/)"
}

variable "services" {
  type        = list(any)
  description = "The list of services allowed to assume this role"
}

variable "policies" {
  type = list(
    object({
      name   = string
      policy = string
    })
  )
  default     = []
  description = "List of inline policies to attach to the role"
}


variable "policy_attachments" {
  type        = list(string)
  default     = []
  description = "List of AWS-managed policy ARNs. Do not use this to attach customer managed account specific policies"
}

variable "tags" {
  description = "The tags for resources"
}


