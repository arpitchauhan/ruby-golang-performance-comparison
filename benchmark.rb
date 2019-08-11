require 'rubygems'
require 'bundler/setup'

Bundler.require(:default)

GO_SOURCE_FILES_NAME = 'program.go'.freeze
GO_EXECUTABLES_NAME = 'program'.freeze
INPUT_FILE = 'input_files/file.csv'.freeze
N = 2

def round(number, digits = 2)
  number.round(digits)
end

def run_command(command, args: [], inform_about_execution: true)
  if inform_about_execution
    puts Rainbow("\nExecuting: \"#{command} #{args.join(" ")}\"").faint
  end

  unless system(command, *args)
    raise "Command execution failed!"
  end
end

def run_command_and_get_time_taken(command, args: [], n: 1)
  times = n.times.map do
    start_time = Time.now
    run_command(command, args: args)
    Time.now - start_time
  end

  times
end

def run_command_and_print_time_taken(command, args: [], n: 1)
  times_taken = run_command_and_get_time_taken(command, args: args, n: n)
  average_time_taken = times_taken.sum.to_f / n
  to_print = Rainbow(
    "\nTimes taken: #{times_taken.map(&method(:round)).join(', ')} seconds. " +
    "Average: #{round(average_time_taken)} seconds"
  ).color(:red)
  puts to_print
end

def execute_ruby_program(path, args: [], n: 1)
  run_command_and_print_time_taken("ruby", args: [path.to_s] + args, n: n)
end

def execute_go_program(directory,
                       source_file = GO_SOURCE_FILES_NAME,
                       args: [],
                       n: 1)
  directory = Pathname.new(directory)
  source_file = directory + source_file
  executable = directory + GO_EXECUTABLES_NAME

  # Build go executable
  run_command("go build -o #{executable} #{source_file}",
              inform_about_execution: false)

  # Run executable
  run_command_and_print_time_taken(executable.to_s, args: args, n: n)
end

puts Rainbow("Running golang programs").yellow

[
  "no_concurrency_read_all_at_once",
  "no_concurrency_read_line_by_line",
  "object_oriented_no_concurrency_read_line_by_line",
  "channel_concurrent",
  "sync_atomic_concurrent",
  ["channel_custom_number_of_goroutines_concurrent", ['50']],
  ["sync_atomic_custom_number_of_goroutines_concurrent", ['50']]
].each do |dirname, args|
  dirname = Pathname.new("golang") + dirname
  arguments = [INPUT_FILE] + (args ? args : [])

  execute_go_program(dirname, args: arguments, n: N)
end

puts Rainbow("\nRunning ruby programs").yellow

[
  "script_single_thread_read_line_by_line.rb",
  "object_oriented_single_thread_read_line_by_line.rb",
  "script_single_thread_read_all_at_once.rb",
  ["script_multiple_threads.rb", ['500']]
].each do |filename, args|
  filename = Pathname.new("ruby") + filename
  arguments = [INPUT_FILE] + (args ? args : [])

  execute_ruby_program(filename, args: arguments, n: N)
end
