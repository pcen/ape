# ape
parse statements as string chunk first
- idea is that expressions that cannot be created within a single c expression are only needed before the statement they are used in
- allows constructing temporary results ahead of statement
- write supplementary code + 'actual' line to source
