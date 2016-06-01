
x = 10
puts "Emitter starting, will emit #{x} events"
x.times do |i|
  puts "metric=emitter-event value=1 type=count"
end
puts "Done emitting metrics"
