require 'csv'

filename = ARGV[0]
number_of_threads = ARGV[1].to_i

raise "Please provide a filename" unless filename

unless ARGV[1].to_i.to_s == ARGV[1] && ARGV[1].to_i >= 1
  raise "Please provide a positive integer as number of threads"
end

Thread.abort_on_exception = true

def allocate_work_among_threads(threads:, work_items:)
  work_items_per_thread = work_items / threads
  leftover_items = work_items % threads

  threads.times.map do |i|
    items = work_items_per_thread
    i < leftover_items ? items + 1 : items
  end
end

def allocate_rows_among_threads(threads:, rows:)
  allocated_work = allocate_work_among_threads(threads: threads, work_items: rows)

  allocated_rows = []

  on = 1

  allocated_work.each do |work_for_thread|
    start_row = on
    end_row = on ? (on + work_for_thread - 1) : nil

    allocated_rows << [start_row, end_row]

    on = if end_row == rows
           nil
         elsif end_row
           end_row + 1
         else
           nil
         end
  end

  allocated_rows
end

def invalid?(row)
  row[0] == '' ||
  row[1] == '' ||
  row[2] == '' ||
  row[3] == '' ||
  row[3].to_i.to_s != row[3].to_s ||
  row[3].to_i < 1 ||
  row[3].to_i > 50
end

data = CSV.read(filename)

number_of_records = data.size

allocated_rows = allocate_rows_among_threads(threads: number_of_threads,
                                             rows: number_of_records)

threads = []

number_of_threads.times do |index|
  start_row, end_row = allocated_rows[index]
  next unless start_row && end_row

  # array indexing starts with zero
  start_row -= 1
  end_row -= 1

  # puts "Thread #{index + 1} processing rows from #{start_row} to #{end_row}"

  threads << Thread.new do
    value = data[start_row..end_row].select do |record|
      invalid?(record)
    end.count

    # puts "Thread #{index + 1} done!"

    value
  end
end

puts "Number of invalid rows = #{threads.map(&:value).sum}"
