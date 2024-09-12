class user_base (
  String $process_name,
) {
  require fastnode_base

  include user_base::executable
  include user_base::fluentd
}
