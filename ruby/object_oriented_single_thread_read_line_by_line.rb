require 'csv'

class CalculateInvalidRowsInCsvFile
  def self.run(filename)
    new(filename).run
  end

  def initialize(filename)
    raise "Please provide a filename" unless filename
    @filename = filename
  end

  def run
    invalid_rows = 0

    each_row do |row|
      invalid_rows += 1 unless Record.valid?(row)
    end

    invalid_rows
  end

  private

  attr_reader :filename

  def each_row
    CSV.foreach(filename) { |row| yield(row) }
  end
end

class Record
  def self.valid?(row)
    new(row).valid?
  end

  def initialize(row)
    @id, @first_name, @last_name, @age = row
  end

  def valid?
    all_present? && age_is_an_integer? && age_is_in_range?
  end

  private

  attr_reader :id, :first_name, :last_name, :age

  def present?(value)
    value != ''
  end

  def all_present?
    present?(id) && present?(first_name) && present?(last_name) && present?(age)
  end

  def age_is_an_integer?
    age.to_i.to_s == age.to_s
  end

  def age_is_in_range?
    age.to_i >= 1 && age.to_i <= 50
  end
end

filename = ARGV[0]

puts "Number of invalid rows = #{CalculateInvalidRowsInCsvFile.run(filename)}"
