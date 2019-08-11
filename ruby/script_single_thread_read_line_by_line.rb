require 'csv'

filename = ARGV[0]
raise "Please provide a filename" unless filename

def invalid?(row)
  row[0] == '' ||
    row[1] == '' ||
    row[2] == '' ||
    row[3] == '' ||
    row[3].to_i.to_s != row[3].to_s ||
    row[3].to_i < 1 ||
    row[3].to_i > 50
end

invalid_rows = 0

# More performant than using select
CSV.foreach(filename) do |row|
  invalid_rows+=1 if invalid?(row)
end

puts "Number of invalid rows = #{invalid_rows}"
